/*
 Copyright 2020 The Magma Authors.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package storage

import (
	"database/sql"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type digestLookup struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

type digestInfo struct {
	network         string
	digest          string
	lastUpdatedTime int64
}

type DigestInfoSet []digestInfo

// GetNetworks returns a list of network IDs for all digests in DigestInfoSet.
func (digestInfoSet DigestInfoSet) GetNetworks() []string {
	ret := []string{}
	for _, digestInfo := range digestInfoSet {
		ret = append(ret, digestInfo.network)
	}
	return ret
}

type DigestLookup interface {
	// Initialize the backing store.
	Initialize() error

	// GetDigests returns a list of digests of that satisfy the filtering criteria
	// specified by the arguments.
	// Caveats:
	// 1. If networks is empty, returns digests for all networks.
	// 2. earliestUpdateTime is recorded in unix seconds. If a positive integer,
	// filter for all digests that were last updated earlier than this time
	// (i.e. outdated digests); otherwise, the time filter is not applied.
	GetDigests(networks []string, earliestUpdateTime int64) (DigestInfoSet, error)

	// SetDigest creates/updates the subscribers digest for a particular network.
	SetDigest(network string, digest string) error

	// DeleteDigests removes digest rows by network IDs.
	DeleteDigests(networks []string) error
}

const (
	digestLookupTableName = "subscriberdb_flat_digests"

	digestLookupNidCol             = "network_id"
	digestLookupDigestCol          = "digest"
	digestLookupLastUpdatedTimeCol = "last_updated_at"
)

func NewDigestLookup(db *sql.DB, builder sqorc.StatementBuilder) DigestLookup {
	return &digestLookup{db: db, builder: builder}
}

func (l *digestLookup) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(digestLookupTableName).
			IfNotExists().
			Column(digestLookupNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(digestLookupDigestCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(digestLookupLastUpdatedTimeCol).Type(sqorc.ColumnTypeBigInt).NotNull().EndColumn().
			PrimaryKey(digestLookupNidCol).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize digest lookup table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *digestLookup) GetDigests(networks []string, earliestUpdateTime int64) (DigestInfoSet, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		filters := squirrel.And{}
		if len(networks) > 0 {
			filters = append(filters, squirrel.Eq{digestLookupNidCol: networks})
		}
		if earliestUpdateTime > 0 {
			filters = append(filters, squirrel.Lt{digestLookupLastUpdatedTimeCol: earliestUpdateTime})
		}

		rows, err := l.builder.
			Select(digestLookupNidCol, digestLookupDigestCol, digestLookupLastUpdatedTimeCol).
			From(digestLookupTableName).
			Where(filters).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "gets digest for networks %+v", networks)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetDigest")

		digestInfoSet := DigestInfoSet{}
		for rows.Next() {
			network, digest, lastUpdatedTime := "", "", int64(0)
			err = rows.Scan(&network, &digest, &lastUpdatedTime)
			if err != nil {
				return nil, errors.Wrap(err, "select digests for networks, SQL row scan error")
			}
			digestInfo := digestInfo{
				network:         network,
				digest:          digest,
				lastUpdatedTime: lastUpdatedTime,
			}
			digestInfoSet = append(digestInfoSet, digestInfo)
		}
		err = rows.Err()
		if err != nil {
			return nil, errors.Wrap(err, "select digests for network, SQL rows error")
		}
		return digestInfoSet, nil
	}

	txRet, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txRet.(DigestInfoSet)
	return ret, nil
}

func (l *digestLookup) SetDigest(network string, digest string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "SetDigest")

		now := clock.Now().Unix()
		_, err := l.builder.
			Insert(digestLookupTableName).
			Columns(digestLookupNidCol, digestLookupDigestCol, digestLookupLastUpdatedTimeCol).
			Values(network, digest, now).
			OnConflict(
				[]sqorc.UpsertValue{
					{Column: digestLookupDigestCol, Value: digest},
					{Column: digestLookupLastUpdatedTimeCol, Value: now},
				},
				digestLookupNidCol,
			).
			RunWith(sc).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "insert digest for network %+v", network)
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *digestLookup) DeleteDigests(networks []string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Delete(digestLookupTableName).
			Where(squirrel.Eq{digestLookupNidCol: networks}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "delete digests")
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func GetDigest(l DigestLookup, network string) (string, int64, error) {
	digestInfoSet, err := l.GetDigests([]string{network}, 0)
	if err != nil {
		return "", 0, err
	}
	// There should be at most 1 digest for each network
	// if digest not found, return default value
	if len(digestInfoSet) == 0 {
		return "", 0, nil
	}
	digestInfo := digestInfoSet[0]
	return digestInfo.digest, digestInfo.lastUpdatedTime, nil
}

func GetOutdatedNetworks(l DigestLookup, updateDeadlineTime int64) ([]string, error) {
	digestInfoSet, err := l.GetDigests([]string{}, updateDeadlineTime)
	if err != nil {
		return nil, err
	}
	networks := digestInfoSet.GetNetworks()
	return networks, nil
}

func GetAllNetworks(l DigestLookup) ([]string, error) {
	digestInfoSet, err := l.GetDigests([]string{}, 0)
	if err != nil {
		return nil, err
	}
	networks := digestInfoSet.GetNetworks()
	return networks, nil
}
