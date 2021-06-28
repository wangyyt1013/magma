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

package subscriberdb_cache_test

import (
	"testing"
	"time"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/lte/cloud/go/services/subscriberdb_cache"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestSubscriberdbCacheWorker(t *testing.T) {
	db, err := test_utils.GetSharedMemoryDB()
	assert.NoError(t, err)
	flatDigestStore := storage.NewFlatDigestLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, flatDigestStore.Initialize())
	perSubDigestStore := storage.NewPerSubDigestLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, perSubDigestStore.Initialize())
	serviceConfig := subscriberdb_cache.Config{
		SleepIntervalSecs:  5,
		UpdateIntervalSecs: 300,
	}

	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	allNetworks, err := storage.GetAllNetworks(flatDigestStore)
	assert.NoError(t, err)
	assert.Equal(t, []string{}, allNetworks)
	flatDigest, err := flatDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	checkDigestEqual(t, "", flatDigest, true)
	perSubDigests, err := perSubDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	_, ok := perSubDigests.(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, map[string]string{}, perSubDigests)

	err = configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	subscriberdb_cache.RenewDigests(flatDigestStore, perSubDigestStore, serviceConfig)
	flatDigest, err = flatDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	checkDigestEqual(t, "", flatDigest, false)
	flatDigestCanon := flatDigest.(storage.DigestInfo).Digest
	perSubDigests, err = perSubDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	_, ok = perSubDigests.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, perSubDigests, "apn")
	assert.NotEqual(t, "", perSubDigests.(map[string]string)["apn"])
	apnDigestCanon := perSubDigests.(map[string]string)["apn"]

	// Detect outdated digests and update
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: lte.APNEntityType, Key: "apn1",
				Config: &lte_models.ApnConfiguration{},
			},
			{
				Type: lte.SubscriberEntityType, Key: "IMSI99999",
				Config: &models.SubscriberConfig{
					Lte: &models.LteSubscription{State: "ACTIVE"},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	clock.SetAndFreezeClock(t, clock.Now().Add(10*time.Minute))
	subscriberdb_cache.RenewDigests(flatDigestStore, perSubDigestStore, serviceConfig)
	flatDigest, err = flatDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	checkDigestEqual(t, flatDigestCanon, flatDigest, false)

	perSubDigests, err = perSubDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	_, ok = perSubDigests.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, perSubDigests, "99999")
	assert.NotEqual(t, "", perSubDigests.(map[string]string)["99999"])
	assert.Equal(t, apnDigestCanon, perSubDigests.(map[string]string)["apn"])
	clock.UnfreezeClock(t)

	// Detect newly added and removed networks
	err = configurator.CreateNetwork(configurator.Network{ID: "n2"}, serdes.Network)
	assert.NoError(t, err)
	configurator.DeleteNetwork("n1")

	clock.SetAndFreezeClock(t, clock.Now().Add(20*time.Minute))
	subscriberdb_cache.RenewDigests(flatDigestStore, perSubDigestStore, serviceConfig)
	flatDigest, err = flatDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	checkDigestEqual(t, "", flatDigest, true)
	perSubDigests, err = perSubDigestStore.GetDigest("n1")
	assert.NoError(t, err)
	_, ok = perSubDigests.(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, map[string]string{}, perSubDigests.(map[string]string))

	flatDigest, err = flatDigestStore.GetDigest("n2")
	assert.NoError(t, err)
	checkDigestEqual(t, "", flatDigest, false)
	perSubDigests, err = perSubDigestStore.GetDigest("n2")
	assert.NoError(t, err)
	_, ok = perSubDigests.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, perSubDigests, "apn")
	assert.NotEqual(t, "", perSubDigests.(map[string]string)["apn"])

	allNetworks, err = storage.GetAllNetworks(flatDigestStore)
	assert.NoError(t, err)
	assert.Equal(t, []string{"n2"}, allNetworks)
	clock.UnfreezeClock(t)
}

func checkDigestEqual(t *testing.T, expected string, digest interface{}, equal bool) {
	digestInfo, ok := digest.(storage.DigestInfo)
	assert.True(t, ok)
	if equal {
		assert.Equal(t, expected, digestInfo.Digest)
	} else {
		assert.NotEqual(t, expected, digestInfo.Digest)
	}
}
