package cloud

type CloudDataBase struct {
	url string
}

func NewCloudDataBase(url string) *CloudDataBase {
	return &CloudDataBase{url}
}

func (db *CloudDataBase) Read() ([]byte, error) {
	return []byte{}, nil
}

func (db *CloudDataBase) Write(content []byte) {

}
