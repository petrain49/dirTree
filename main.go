package main


func main() {
	makeTree(".")
}

func makeTree(pathToFile string) error {
	files, err := getContent(pathToFile)
	if err != nil {
		return err
	}

	var depth uint

	err = printFiles(files, depth)
	if err != nil {
		return err
	}
	return nil
}