package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	. "github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
func processFile(srcFile string) {
	f, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tarReader := tar.NewReader(gzf)

	i := 0
	buf := make([]byte, 1024*1024)
	buf2 := make([]byte, 1024*1024)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if strings.HasSuffix(name, "_icgraph.dot") {
				//fmt.Println("(", i, ")", "Name: ", name)
				len, err := tarReader.Read(buf)
				fbuf := buf[:len]
				if err == nil {
					for {
						len2, err := tarReader.Read(buf2)
						fbuf = BytesCombine(fbuf, buf2[:len2])
						if err != nil {
							break
						}
					}
				}

				ParseICallGraphByBuffer(fbuf)
			}

		default:
			fmt.Printf("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}

		i++
	}
}

var tarCmd *cobra.Command = &cobra.Command{
	Use:   "tar",
	Short: "full collaboration graph from tar file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()

		//filter := cmd.Flag("filter").Value.String()
		ParseICallGraphStart()
		processFile(source)
		result := ParseICallGraphEnd()
		if cmd.Flag("findCrossRefs").Value.String() == "true" {
			crossRefs := result.FindCrossRef(MergeHeaderFunc)
			for _, cf := range crossRefs {
				fmt.Println(cf)
			}
			return
		}
		if cmd.Flag("mergePackage").Value.String() == "true" {
			result = result.MergeHeaderFile(MergePackageFunc)
		}
		if cmd.Flag("java").Value.String() != "" && cmd.Flag("common").Value.String() != "" {
			javaFilterFile := cmd.Flag("java").Value.String()
			commonFilterFile := cmd.Flag("common").Value.String()
			javaFilter := CreatePrefixFilter(javaFilterFile)

			commonFilter := CreatePrefixFilter(commonFilterFile)
			printRelation(result, javaFilter.Match, commonFilter.Match)
			return

		}

		result.ToDot(cmd.Flag("output").Value.String(), ".", func(s string) bool {
			return false
		})
	},
}

func init() {
	rootCmd.AddCommand(tarCmd)
	tarCmd.Flags().BoolP("findCrossRefs", "C", false, "find cross references")
	tarCmd.Flags().BoolP("mergePackage", "P", false, "merge package/folder for include dependencies")
	tarCmd.Flags().StringP("source", "s", "", "source code directory")
	tarCmd.Flags().StringP("filter", "f", "coll__graph.dot", "dot file filter")
	tarCmd.Flags().StringP("output", "o", "dep.dot", "output dot file name")
	tarCmd.Flags().StringP("java", "j", "", "java class filter")
	tarCmd.Flags().StringP("common", "c", "", "common java class")
}
