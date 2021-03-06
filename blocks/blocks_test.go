package blocks

import (
	"fmt"
	. "github.com/Inflatablewoman/blocker/gocheck2"
	. "gopkg.in/check.v1"
	"os"
	"testing"
	"time"
)

func Test(t *testing.T) {
	TestingT(t)
}

type BlockSuite struct {
}

var _ = Suite(&BlockSuite{})

const inputFile = "testdata/tempest.txt"
const bibleInFile = "testdata/kjv.txt"
const bibleOutFileName = "kjv.txt"
const changedInputFile = "testdata/tempest_changed.txt"
const changedAgainInputFile = "/testdata/tempest_changed_again.txt"
const outputFileName = "tempest.txt"
const changedOutputFileName = "tempest_changed.txt"
const liteIdeInFile = "testdata/liteidex23.2.linux-32.tar.bz2"
const liteIdeoutFile = "liteidex23.2.linux-32.tar.bz2"

func (s *BlockSuite) SetUpSuite(c *C) {
	// Ensure string is to lower
	StorageProviderName = "nfs"

	// Now set up repos
	SetUpRepositories()
}

func (s *BlockSuite) TestAdvancedBlockCopyDelete(c *C) {

	// Block a file...
	start := time.Now()
	bibleBlockFile, err := BlockFile(liteIdeInFile)
	end := time.Now()

	fmt.Printf("1st Blocked LiteIde took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	bibleOutFile := os.TempDir() + "/" + liteIdeoutFile

	// Clean up any old file
	os.Remove(bibleOutFile)

	// Get the file and create a copy to the output
	start = time.Now()
	err = UnblockFile(bibleBlockFile.ID, bibleOutFile)
	end = time.Now()

	fmt.Printf("1st Unblocked LiteIde took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// Maintain the first BlockFileID
	firstBlockFileID := bibleBlockFile.ID
	firstBlockFileBlockHash := bibleBlockFile.BlockList[0].Hash

	// Check the block store has the data...
	blockExists, _ := BlockStore.CheckBlockExists(firstBlockFileBlockHash)
	c.Assert(blockExists, IsTrue)

	// Block the file again
	start = time.Now()
	bibleBlockFile, err = BlockFile(liteIdeInFile)
	end = time.Now()

	fmt.Printf("2nd Blocked LiteIde took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// Maintain the second BlockFileID
	secondBlockFileID := bibleBlockFile.ID
	secondBlockFileBlockHash := bibleBlockFile.BlockList[0].Hash

	// Should be new BlockedFileID despite the same content
	c.Assert(secondBlockFileID != firstBlockFileID, IsTrue)

	// Check that block used in first block is the same
	c.Assert(firstBlockFileBlockHash == secondBlockFileBlockHash, IsTrue)

	// Copy our first BlockedFile
	start = time.Now()
	bibleBlockFile, err = CopyBlockedFile(firstBlockFileID)
	end = time.Now()

	fmt.Printf("Copied BlockFile: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	thirdBlockFileID := bibleBlockFile.ID
	thirdBlockFileBlockHash := bibleBlockFile.BlockList[0].Hash

	// Should be new BlockedFileID despite the same content
	c.Assert(thirdBlockFileID != firstBlockFileID, IsTrue)

	// Check that block used in first block is the same
	c.Assert(firstBlockFileBlockHash == thirdBlockFileBlockHash, IsTrue)

	fileBlockInfo, err := FileBlockInfoStore.GetFileBlockInfo(thirdBlockFileBlockHash)

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// This file block should be used 3 times
	c.Assert(fileBlockInfo.UseCount == 3, IsTrue, Commentf("Block shound be used 3 time.  Is used: : %v", fileBlockInfo.UseCount))

	// Delete first block...
	start = time.Now()
	err = DeleteBlockedFile(firstBlockFileID)
	end = time.Now()

	fmt.Printf("Delete Block took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// Ensure the BlockedFile is no longer there
	_, err = BlockedFileStore.GetBlockedFile(firstBlockFileID)

	// We should have an error
	c.Assert(err == nil, IsFalse)

	// Check the use count
	fileBlockInfo, err = FileBlockInfoStore.GetFileBlockInfo(thirdBlockFileBlockHash)

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// This file block should be used 2 times
	c.Assert(fileBlockInfo.UseCount == 2, IsTrue, Commentf("Block shound be used 2 time.  Is used: : %v", fileBlockInfo.UseCount))

	// Delete second blockedfile
	start = time.Now()
	err = DeleteBlockedFile(secondBlockFileID)
	end = time.Now()

	fmt.Printf("Delete Block took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// Delete third blockedfile
	start = time.Now()
	err = DeleteBlockedFile(thirdBlockFileID)
	end = time.Now()

	fmt.Printf("Delete Block took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// There should now be no reference to the data in any repository

	// Check the use count
	_, err = FileBlockInfoStore.GetFileBlockInfo(thirdBlockFileBlockHash)

	// There should be an error
	c.Assert(err == nil, IsFalse)

	// Check the block store has deleted the data...
	blockExists, _ = BlockStore.CheckBlockExists(thirdBlockFileBlockHash)
	c.Assert(blockExists, IsFalse)
}

func (s *BlockSuite) TestLiteide(c *C) {

	//c.Skip("Want to work out what is going on")

	/*BlockSize = BlockSize4Mb
	UseCompression = true
	UseEncryption = true*/

	// Block the bigger
	start := time.Now()
	bibleBlockFile, err := BlockFile(liteIdeInFile)
	end := time.Now()

	fmt.Printf("Blocked LiteIde took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	bibleOutFile := os.TempDir() + "/" + liteIdeoutFile

	// Clean up any old file
	os.Remove(bibleOutFile)

	// Get the file and create a copy to the output
	start = time.Now()
	err = UnblockFile(bibleBlockFile.ID, bibleOutFile)
	end = time.Now()

	fmt.Printf("Unblocked LiteIde took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	err = DeleteBlockedFile(bibleBlockFile.ID)
	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))
}

func (s *BlockSuite) TestKingJamesBible(c *C) {

	//c.Skip("Want to work out what is going on")

	/*BlockSize = BlockSize4Mb
	UseCompression = true
	UseEncryption = true*/

	// Block the bigger
	start := time.Now()
	bibleBlockFile, err := BlockFile(bibleInFile)
	end := time.Now()

	fmt.Printf("Blocked King James Bible took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	bibleOutFile := os.TempDir() + "/" + bibleOutFileName

	// Clean up any old file
	os.Remove(bibleOutFile)

	// Get the file and create a copy to the output
	start = time.Now()
	err = UnblockFile(bibleBlockFile.ID, bibleOutFile)
	end = time.Now()

	fmt.Printf("Unblocked King James Bible took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	err = DeleteBlockedFile(bibleBlockFile.ID)
	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))
}

func (s *BlockSuite) TestTempest(c *C) {

	/*BlockSize = BlockSize4Mb
	UseCompression = true
	UseEncryption = true*/

	// Get some info about the file we are going test
	inputFileInfo, _ := os.Stat(inputFile)

	// Block the file
	start := time.Now()
	blockFile, err := BlockFile(inputFile)
	end := time.Now()

	fmt.Printf("Blocked Tempest took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// Check we have an ID
	c.Assert(blockFile.ID != "", IsTrue)

	// Check we read the full file size
	c.Assert(blockFile.Length == inputFileInfo.Size(), IsTrue)

	// Make sure the item returned some blocks
	c.Assert(len(blockFile.BlockList) > 0, IsTrue)

	outputFile := os.TempDir() + "/" + outputFileName

	// Clean up any old file
	os.Remove(outputFile)

	// Get the file and create a copy to the output
	start = time.Now()
	err = UnblockFile(blockFile.ID, outputFile)
	end = time.Now()

	fmt.Printf("Unblocked Tempest took: %v\n", end.Sub(start))

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// Get some info about the file we are going test
	outputFileInfo, _ := os.Stat(outputFile)

	// Check we wrote the full file size
	c.Assert(outputFileInfo.Size() == inputFileInfo.Size(), IsTrue)

	err = DeleteBlockedFile(blockFile.ID)
	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))
}

// Benchmark the blocking of Tempest file
/*func (s *BlockSuite) BenchmarkTempestCompressedEncrypted30Kb(c *C) {
	for i := 0; i < c.N; i++ {
		// Set up test
		BlockSize = BlockSize30Kb
		UseCompression = true
		UseEncryption = true
		BlockFile(inputFile)
	}
}

func (s *BlockSuite) BenchmarkTempestCompressedEncrypted4Mb(c *C) {
	for i := 0; i < c.N; i++ {
		// Set up test
		BlockSize = BlockSize4Mb
		UseCompression = true
		UseEncryption = true
		BlockFile(inputFile)
	}
}

func (s *BlockSuite) BenchmarkTempestUncompressedUnencrypted4Mb(c *C) {
	for i := 0; i < c.N; i++ {
		// Set up test
		BlockSize = BlockSize4Mb
		UseCompression = false
		UseEncryption = false
		// Need to clean out the block store directory

		BlockFile(inputFile)
	}
}*/

func (s *BlockSuite) TestChangeTempest(c *C) {

	c.Skip("Not what I want to test right now")

	// Set up test
	BlockSize = BlockSize30Kb
	UseCompression = true
	UseEncryption = true

	// Get some info about the file we are going test
	changedInputFileInfo, _ := os.Stat(changedInputFile)

	blockFile, err := BlockFile(inputFile)

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	firstFileHash := blockFile.BlockList[0].Hash

	// Block the file again.
	blockFile, err = BlockFile(inputFile)

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// Check we have an ID
	c.Assert(blockFile.ID != "", IsTrue)

	// Remember the first block id
	firstBlockID := blockFile.ID

	// Check that block used in first block is the same
	c.Assert(firstFileHash == blockFile.BlockList[0].Hash, IsTrue)

	// Block the file again.  New version should be created
	blockFile, err = BlockFile(changedInputFile)

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// Check we have an ID
	c.Assert(blockFile.ID != "", IsTrue)

	// Remember the first block id
	secondBlockID := blockFile.ID

	changedOutputFile := os.TempDir() + "/" + changedOutputFileName

	// Clean up any old file
	os.Remove(changedOutputFile)

	// Get the file and create a copy to the output
	err = UnblockFile(blockFile.ID, changedOutputFile)

	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	// Get some info about the file we are going test
	outputFileInfo, _ := os.Stat(changedOutputFile)

	// Check we wrote the full file size
	c.Assert(outputFileInfo.Size() == changedInputFileInfo.Size(), IsTrue)

	err = DeleteBlockedFile(firstBlockID)
	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))

	err = DeleteBlockedFile(secondBlockID)
	// No error
	c.Assert(err == nil, IsTrue, Commentf("Failed with error: %v", err))
}
