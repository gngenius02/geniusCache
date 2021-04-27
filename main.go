package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

var (
	wgMain sync.WaitGroup
	wg     sync.WaitGroup
)

func checkDBForHash(seed string, db *badger.DB) {
	defer wgMain.Done()
	hashit := func(s string) string {
		dig := sha256.Sum256([]byte(s))
		return hex.EncodeToString(dig[:])
	}

	wg.Add(1)
	go db.View(exists(seed, seed))

	hash := hashit(seed)
	for i := 0; i <= 1000000; i++ {
		wg.Add(1)
		go db.View(exists(hash, seed))
		hash = hashit(hash)
		if i == 1000000 {
			saveToDb(seed, hash, db)
		}
	}
	wg.Wait()
}

func main() {
	runtime.GOMAXPROCS(128)
	start := time.Now()

	db := connectToDB()
	defer db.Close()

	for i := 0; i < 48; i++ {
		rand := getRandBytes()
		if i%12 == 0 {
			fmt.Println(rand)
		}
		wgMain.Add(1)
		go checkDBForHash(rand, db)
	}
	wgMain.Wait()

	fmt.Println("finished 48 random checks in", time.Since(start))

	// db.Flatten(32)
	// checkDBForHash("422f9fd47ef453b532644662253b71ee52e9b5ed175e656faeee040ad5afc9e7", db)
	// checkDBForHash("fc569185b84499cde2231eb352faeab49d7b2e86375fd984140bb581c104560a", db)
	// checkDBForHash("2cc7f343907def3023f6ac87084428a6169d6972dab5b4827eb310f29dcbd222", db)

	//	getFileList("/home/node/oldstuff", db)		//Load data into DB.
	//	checkDBForHash(os.Args[1], db)
	//	testhashes := []string{"3f7f9b47cc4d4efe3bfcfd5e01c2eacec61df23c048a93a9d98f729c440a3148", "061eefc3d88536ef377ecbb01f436af539f689d7b35b81bb8591cfe3a8d40e5e", "3f67a354dc002a3caa8d7e6c6b09d69c4a8920f977256eb204f55623fa4d9ed8", "d7d40f9def6b0d3b39ea6c301a9228e560e145966b463907fa9f0514ba62bdcb", "6795bb0b9eef3c4898cd173806577743b90670cf15e761e89ed7bff26906cdc5", "dcf12ee1a52f7e6753569f86c94f8ab685c11b879773963bbdd7235391604e1c", "dcfd798c156933f4474a9fc82e3a40df2e800b359bb756a8bebf5e278d9e6617", "8b1e85c18a229793423cd03f1903acbe03b02ac2b1281e205ce29876f4be350f", "b3be5ce98646c499584a1fa6d2fb09aaa9338e02152efd5cf2d3c000b23569fd", "5063019f378fbdba0acfab23117dc5d518aacc0b4d79f8c984ddd655be37934e", "f353c38c61165f4f346f9424924434a7a2f656d47fa9f6d43dad64ac3cfdd294", "47b3a82ef5b3da3abffa982dc086172ecf302af7c84af4dec71592e167d57f85", "1d6cb7abb433605e063f353d17de8340072a546e32b83305992fc7bbcee774ea", "9e03555ef4d73be520475f250ed921c794569e3ff78c2e9c712e14dc564489b2", "bdddf894c4d1b02268ac3b023d555a4c13c5792e696a40b013f77d5fe90c7310", "da0485e8ed45a2dab2d945159ed895a7fa4df1d19a01f98f7197fff5183949b9", "82de4791d8a1e41879ed43ad5009369e94ac0bdfcc4d83b4c10fd97ecbbfcb67", "e2d06b00bcb4a4b96fd481a7fed363311011da4aa183d924f782917a89bbffcc", "0aecfe9bb7ba01f3667287cf4b51142486b7ce4c99a292d853e6e5e126b70073", "6a8a63d09842131ae8e88e469de06db6e40c0a18b9d66b22204468049fc9cfa8", "b5416a69bb683f4ba66cba029a061ee38da40b25922b6f984afdac493a9c169d", "a1df6b851e6a93828fe14b60d24ea76b73f4292f68f6d65dc2390e8e22ae1850", "1cb16cb7ca5a4c64992877ec36bcd9324e9adb5775b573396e1abbb2cdf6c69a", "02fb3cda569bd8e29ff1d90bbb6c051f66f24f485256c77b3dee336e192df38d", "cc215c1c0bddd48e7e9c6831f689293be407df48e291e96553e6334a8e8f8728", "22c87326ea4889f731b5b537eba6d256cbe95f4e40111466e50e49d552947873", "e9b616374f938c28bf10455a7e57417d8a385c04c8dbb350fc2ce925a7967cab", "cdee24525c1702a4412c2c7195e300768af6b052e7ff1a10975f1be44b30ce8f", "aced7b4dfd3831839d86bef6b048bffdf454807146df7847e08e8f3066f5dba4", "d9346e81700cac85337ad3dece0c4cec1de6b7b99a2f6129d87c3bcd6b381b0d", "deaaa636cfbac29d96ebeb0ca7abd97f850bee703dcec575c4c9d3e170d38034", "10a96cf653a2c5f7906598c24926d3e174c6d263e8fb952bff0a0d0c66aab1e9", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}
	//	for _, v := range testhashes {
	//		checkDBForHash(v, db)
	//	}
}
