package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

var (
	version = "0.0.0"
	commit  = ""
)

func ExtractTarGz(gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		mode := header.FileInfo().Mode()

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, mode); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(header.Name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(f, tarReader); err != nil {
				return err
			}

		default:
			err := fmt.Errorf("tar: unknown type %q in %q", header.Typeflag, header.Name)
			return err
		}

	}
	return nil
}

func main() {
	log.SetFlags(0)

	exe := path.Base(os.Args[0])
	if len(os.Args) != 2 {
		log.Fatalln("Usage:", exe, "<input_tgz_file>")
	}

	if os.Args[1] == "--version" {
		fmt.Println(exe, version, commit)
		os.Exit(0)
	}

	r, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Close()

	if err := ExtractTarGz(r); err != nil {
		log.Fatalln(err)
	}
}
