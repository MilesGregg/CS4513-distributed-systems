package replication

import (
    "errors"
    "fmt"
    "os"
    "time"
    "strconv"
)
/*
TODO: implement the methods at the bottom of this file.
 */

/*
InitializeNodes creates folders for each "node" computer in the mini-S3 system.

File directory will be /src/replication/nodes/NODE_NUMBERS/BUCKET_NAME/FILENAME
 */
func InitializeNodes(nodes int) {
    setNumberNodes(nodes)
    for i := 0; i < nodes; i++ {
        nodeDir := fmt.Sprintf("nodes/%d", i)
        err := os.MkdirAll(nodeDir, os.ModePerm)
        checkError(err)
    }
}

/*
ResetNodes removes /nodes/ and all subdirectories and all files within those directories
 */
func ResetNodes() {
    err := os.RemoveAll("nodes")
    checkError(err)
}



/*
BucketExists determines if the specified bucket exists
*/
func BucketExists(bucketName string) bool {
    numberNodes := getNumberNodes()
    for i := 0; i < numberNodes; i++ {
        bucket := fmt.Sprintf("nodes/%d/%s", i, bucketName)
        if _, err := os.Stat(bucket); errors.Is(err, os.ErrNotExist) {
            // bucket does not exist at node i
            return false
        }
    }
    return true
}

/*
CreateBucket should create a fake S3 bucket in the form of a directory. A bucket must be created before a file can be written
to the bucket. The bucket must be created in each node.
*/
func CreateBucket(
    bucketName string,
) {
    // TODO: implement this method
    for i := 0; i < getNumberNodes(); i++ {
        err := os.MkdirAll("nodes/" + strconv.Itoa(i) + "/" + bucketName, os.ModePerm)
        checkError(err)
    }
}


/*
WriteNodeFile should write a byte array (contents) to the specified bucket to the specified node (nodeIndex) with the
specified file name.

Should write the file version to a paired file (perhaps something like fileName.version?). version should be the current
time.

Returns the number of bytes written
 */
func WriteNodeFile(
    nodeIndex int,
    bucketName string,
    fileName string,
    contents []byte,
    version time.Time,
) int {
    // TODO: implement this method.
    err := os.MkdirAll("nodes/" + strconv.Itoa(nodeIndex) + "/" + bucketName, os.ModePerm)
    checkError(err)

    // check to see if the file exists first before creating it
	filename_absolute_path := "nodes/" + strconv.Itoa(nodeIndex) + "/" + bucketName + "/" + fileName
    check_file, err := os.Open(filename_absolute_path)
    defer check_file.Close()
    if err != nil {
        // file does not exist, so create it
        file, err := os.Create(filename_absolute_path)
        checkError(err)
        defer file.Close() // close at the very end
    }

    // write the incoming contents into the file
    err = os.WriteFile(filename_absolute_path, contents, os.ModePerm)
    checkError(err)
    
    // create version file
    version_file, err := os.Create(filename_absolute_path + ".version")
    checkError(err)
    defer version_file.Close()

    // write into the version file
    _, err = version_file.Write([]byte(version.Format(time.RFC3339)))
    checkError(err)

    return len(contents)
}

/*
ReadNodeFile should read the specified file from the specified node from the specified bucket

Returns file contents, file version
 */
func ReadNodeFile(
    nodeIndex int,
    bucketName string,
    fileName string,
) ([]byte, time.Time) {
    // TODO: implement this method.
	
	// absolute filepath to the current file we are working on 
	filename_absolute_path := "nodes/" + strconv.Itoa(nodeIndex) + "/" + bucketName + "/" + fileName
	
	// open up the file
    file, error_file := os.Open(filename_absolute_path)
	checkError(error_file)
	defer file.Close()

	// open up the version file
    file_version, error_version := os.Open(filename_absolute_path + ".version")
    checkError(error_version)
    defer file_version.Close()

    file_metadata, err := file.Stat()
    checkError(err)

    file_information := make([]byte, file_metadata.Size())
    _, err = file.Read(file_information)
    checkError(err)

    var info []byte
    _, err = fmt.Fscan(file_version, &info)
    checkError(err)
    version, err := time.Parse(time.RFC3339, string(info))
    checkError(err)

    return file_information, version

    //return []byte("TODO"), time.Now()
}


