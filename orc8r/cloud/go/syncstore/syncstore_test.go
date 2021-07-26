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

package syncstore_test

import (
	"testing"
	"time"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/clock"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestSyncStore(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewSQLBlobStorageFactory("last_resync_time", db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	store := syncstore.NewSyncStore(db, sqorc.GetSqlBuilder(), fact)
	assert.NoError(t, store.Initialize())

	expectedDigestTree := &protos.DigestTree{
		RootDigest: &protos.Digest{Md5Base64Digest: "root_digest_apple2"},
		LeafDigests: []*protos.LeafDigest{
			{Id: "2", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_banana2"}},
			{Id: "3", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_cherry2"}},
			{Id: "4", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_dragonfruit"}},
		},
	}
	expectedDigestTree2 := &protos.DigestTree{
		RootDigest: &protos.Digest{Md5Base64Digest: "root_digest_banana"},
	}
	objs1 := map[string][]byte{
		"1": []byte("apple"),
	}
	objs2 := map[string][]byte{
		"2": []byte("banana"),
		"3": []byte("cherry"),
	}
	expectedObjs := [][]byte{
		[]byte("apple"),
		[]byte("banana"),
		[]byte("cherry"),
	}

	t.Run("initially empty", func(t *testing.T) {
		digestTrees, err := store.GetDigests([]string{"n0", "n1"}, time.Now().Unix(), true)
		assert.NoError(t, err)
		assert.Empty(t, digestTrees)

		page, nextToken, err := store.GetCachedByPage("n0", "", 10)
		assert.NoError(t, err)
		assert.Empty(t, page)
		assert.Empty(t, nextToken)
	})

	t.Run("basic insert digests", func(t *testing.T) {
		expectedDigestTree := &protos.DigestTree{
			RootDigest: &protos.Digest{Md5Base64Digest: "root_digest_apple"},
			LeafDigests: []*protos.LeafDigest{
				{Id: "1", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_apple"}},
				{Id: "2", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_banana"}},
				{Id: "3", Digest: &protos.Digest{Md5Base64Digest: "leaf_digest_cherry"}},
			},
		}
		err := store.SetDigest("n0", expectedDigestTree)
		assert.NoError(t, err)

		digestTrees, err := store.GetDigests([]string{"n0"}, time.Now().Unix(), true)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.True(t, proto.Equal(expectedDigestTree, digestTrees["n0"]))

		digestTrees, err = store.GetDigests([]string{"n0"}, time.Now().Unix(), false)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.Equal(t, "root_digest_apple", digestTrees["n0"].RootDigest.Md5Base64Digest)
		assert.Empty(t, digestTrees["n0"].GetLeafDigests())
	})

	t.Run("upsert digests", func(t *testing.T) {
		err = store.SetDigest("n0", expectedDigestTree)
		assert.NoError(t, err)

		digestTrees, err := store.GetDigests([]string{"n0"}, time.Now().Unix(), true)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.True(t, proto.Equal(expectedDigestTree, digestTrees["n0"]))
	})

	t.Run("get outdated digests", func(t *testing.T) {
		clock.SetAndFreezeClock(t, clock.Now().Add(200*time.Second))
		err = store.SetDigest("n1", expectedDigestTree2)
		assert.NoError(t, err)

		digestTrees, err := store.GetDigests([]string{"n0", "n1"}, clock.Now().Unix(), true)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.True(t, proto.Equal(expectedDigestTree, digestTrees["n0"]))
		assert.Contains(t, digestTrees, "n1")
		assert.True(t, proto.Equal(expectedDigestTree2, digestTrees["n1"]))

		digestTrees, err = store.GetDigests([]string{"n0", "n1"}, clock.Now().Unix()-100, true)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.NotContains(t, digestTrees, "n1")

		digestTrees, err = store.GetDigests([]string{"n0", "n1"}, clock.Now().Unix()-300, true)
		assert.NoError(t, err)
		assert.NotContains(t, digestTrees, "n0")
		assert.NotContains(t, digestTrees, "n1")

		clock.UnfreezeClock(t)
	})

	t.Run("basic insert and get from cache", func(t *testing.T) {
		writer, err := store.UpdateCache("n0")
		assert.NoError(t, err)

		err = writer.InsertMany(objs1)
		assert.NoError(t, err)
		err = writer.InsertMany(objs2)
		assert.NoError(t, err)
		err = writer.Apply()
		assert.NoError(t, err)

		objs, err := store.GetCachedByID("n0", []string{"1", "2", "3"})
		assert.NoError(t, err)
		assert.Equal(t, expectedObjs, objs)

		expectedNextToken, err := configurator_storage.SerializePageToken(&configurator_storage.EntityPageToken{
			LastIncludedEntity: "3",
		})
		assert.NoError(t, err)
		objs, nextToken, err := store.GetCachedByPage("n0", "", 3)
		assert.NoError(t, err)
		assert.Equal(t, expectedObjs, objs)
		assert.Equal(t, expectedNextToken, nextToken)

		objs, nextToken, err = store.GetCachedByPage("n0", nextToken, 3)
		assert.NoError(t, err)
		assert.Empty(t, objs)
		assert.Empty(t, nextToken)

		// When the changes have been applied, the cache writer could no longer be used for insertions
		err = writer.InsertMany(objs1)
		assert.EqualError(t, err, "attempt to insert into network n0 with invalid cache writer")
	})

	t.Run("garbage collection", func(t *testing.T) {
		err = store.SetDigest("n0", expectedDigestTree)
		assert.NoError(t, err)
		err = store.SetDigest("n1", expectedDigestTree2)
		assert.NoError(t, err)
		writer, err := store.UpdateCache("n0")
		assert.NoError(t, err)
		err = writer.InsertMany(objs1)
		assert.NoError(t, err)
		err = writer.InsertMany(objs2)
		assert.NoError(t, err)
		err = writer.Apply()
		assert.NoError(t, err)

		digestTrees, err := store.GetDigests([]string{}, clock.Now().Unix(), true)
		assert.NoError(t, err)
		assert.Contains(t, digestTrees, "n0")
		assert.Contains(t, digestTrees, "n1")
		objs, _, err := store.GetCachedByPage("n0", "", 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, objs)

		// Only track data from network n1
		err = store.CollectGarbage([]string{"n1"})
		assert.NoError(t, err)

		digestTrees, err = store.GetDigests([]string{}, clock.Now().Unix(), true)
		assert.NoError(t, err)
		assert.NotContains(t, digestTrees, "n0")
		assert.Contains(t, digestTrees, "n1")
		objs, _, err = store.GetCachedByPage("n0", "", 10)
		assert.NoError(t, err)
		assert.Empty(t, objs)
	})
}
