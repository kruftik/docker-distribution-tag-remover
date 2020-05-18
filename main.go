package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"
)

var (
	RepoRoot  = flag.String("repository-root", "", "path to docker-distribution storage root")
	ImageLabel  = flag.String("image", "", "image label to use (e.g. minio/mc)")
	ImageTag  = flag.String("tag", "", "image tag to delete (e.g. latest)")
)

type BlankStruct struct {}

func init () {
	flag.Parse()

	if *RepoRoot == "" {
		log.Error("'repository-root' flag is not provided")
	}
	if *ImageLabel == "" {
		log.Error("'image' flag is not provided")
	}
	//if *ImageTag == "" {
	//	log.Error("'tag' flag is not provided")
	//}
	//
	//*ImageTag = "release-26-0-0-361"
}

func main (){
	if ok := IsDir(*RepoRoot + "/docker/registry/v2/repositories"); !ok {
		log.Error("'repository-root' is not valid: it does not contain docker/registry/v2/repositories directory inside")
		return
	}



	r, err := NewRegistry(*RepoRoot, *ImageLabel)
	if err != nil {
		log.Error("cannot init registry: %s", err)
		return
	}

	//tagManifest, err := registry.GetManifestForTag(*ImageTag)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//log.Info("tag '%s' has manifest sha256:%s", *ImageTag, tagManifest)
	//
	//tagsForManifest, err := registry.GetTagsWithManifest(tagManifest)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	//tagsNum := len(tagsForManifest)
	//
	//log.Info("tag '%s' has manifest 'sha256:%s'", *ImageTag, tagManifest)
	//
	//if tagsNum > 1 {
	//	for _, otherTag := range tagsForManifest {
	//		if otherTag != *ImageTag {
	//			log.Info("tag '%s' has the same manifest 'sha256:%s' as %s", otherTag, tagManifest, *ImageTag)
	//		}
	//	}
	//}

	//isSaveToRemove, err := registry.IsTagSaveToRemove(*ImageTag)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	//if !isSaveToRemove {
	//	log.Warn("probably there are tags with the same manifest as %s, it can't be removed", *ImageTag)
	//
	//	otherTags, err := registry.GetTagsWithSameManifest(*ImageTag, false)
	//	if err != nil {
	//		log.Error(err)
	//		return
	//	}
	//
	//	log.Info(otherTags)
	//}

	if flag.NArg() == 0 {
		log.Error("tags-file not provided")
		return
	}
	tagsFile := flag.Arg(0)
	if !IsFile(tagsFile) {
		log.Error("tags file %s does not exist", tagsFile)
	}

	f, err := os.Open(tagsFile)
	if err != nil {
		log.Error("cannot open tags file %s: %s", tagsFile, err)
	}
	buf := bufio.NewReader(f)

	tagsList := make(map[string]struct{})

	for {
		tag, err := buf.ReadString('\n')
		if err == io.EOF {
			log.Info("tags file have been read")
			break
		}

		if err != nil {
			log.Error("cannot read tags file: %s", err)
			return
		}

		tag = strings.TrimRight(tag, "\n")
		tagsList[tag] = BlankStruct{}
	}

	ableToDelete := func(inTag string) bool {
		for tag, _ := range tagsList {
			if tag == inTag {
				return true
			}
		}
		return false
	}

	for tagC, _ := range tagsList {
		sameManifestTags, err := r.GetTagsWithSameManifest(tagC, false)
		if err != nil {
			log.Error(err)
			return
		}
		sameManifestTagsCount := len(sameManifestTags)

		if sameManifestTagsCount == 1 {
			continue
		}

		//allTagsAreSafeToDelete := true
		for _, sameTag := range sameManifestTags {
			if sameTag == tagC {
				continue
			}

			if !ableToDelete(sameTag){
				log.Info("[%s] %s has the same manifest and isn't allowed to be deleted", tagC, sameTag)
				delete(tagsList, tagC)
				//allTagsAreSafeToDelete = false
				break
			}
		}

		//if !allTagsAreSafeToDelete {
		//	log.Info("[%s] there are other tags with the same manifest", tagC)
		//}
	}

	for tag := range tagsList {
		log.Info("%s - delete OK", tag)
	}
}


