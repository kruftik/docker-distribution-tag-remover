package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	REPO_IMAGES_DIR = "docker/registry/v2/repositories"
	IMAGE_TAGS_DIR_POSTFIX = "_manifests/tags"
	CURRENT_TAG_REVISTION_HASH_POSTFIX = "current/link"
)

type TagToManifestMap map[string]string
type ManifestToTagMap map[string][]string

type Registry struct {
	imagesDir string
	imageLabelDir string

	t2m TagToManifestMap
	m2t ManifestToTagMap
}

func NewRegistry(repoRoot, image string) (*Registry, error) {
	r := &Registry{
		t2m: make(TagToManifestMap),
		m2t: make(ManifestToTagMap),
	}

	err := r.init(repoRoot, image)
	if err != nil {
		return nil, fmt.Errorf("cannot create registry: %w", err)
	}

	return r, nil
}

func (r *Registry) init(repoRoot, image string) error {
	r.imagesDir = ImagesDir(repoRoot)

	if !IsDir(r.imagesDir) {
		return fmt.Errorf("%s does not exist", r.imagesDir)
	}
	log.Debug("Using image dir: %s", r.imagesDir)

	r.imageLabelDir = ImageLabelDir(image)

	if !IsDir(r.imageLabelDir) {
		return fmt.Errorf("%s does not exist", r.imageLabelDir)
	}

	return nil
}

func (r *Registry) buildTagAndManifestMaps() error {
	if len(r.t2m) > 0 {
		log.Debug("buildTagAndManifestMaps: maps already built")
		return nil
	}

	getManifestForTag := func(tagDir string) (string, error) {
		tag := filepath.Base(tagDir)

		hashFile := tagDir + "/" + CURRENT_TAG_REVISTION_HASH_POSTFIX

		tagHash, err := ioutil.ReadFile(hashFile)
		if err != nil {
			return "", fmt.Errorf("cannot read %s: %w", hashFile, err)
		}

		tagHashStr := (strings.Split(string(tagHash), ":"))[1]

		log.Debug("Tag %s has manifest hash: %s", tag, tagHashStr)

		return tagHashStr, nil
	}

	addTagToTagsList := func(r *Registry, manifest, tag string) {
		if _, ok := r.m2t[manifest]; !ok {
			r.m2t[manifest] = make([]string, 0)
		}

		r.m2t[manifest] = append(r.m2t[manifest], tag)
	}

	tagsDir := r.imageLabelDir + "/" + IMAGE_TAGS_DIR_POSTFIX + "/*"
	log.Debug("Scanning " + tagsDir)

	tags, err := filepath.Glob(tagsDir)
	if err != nil {
		return fmt.Errorf("cannot list tags for image: %w", err)
	}

	for _, tagDir := range tags {
		tag := filepath.Base(tagDir)
		log.Debug("tag dir: %s", tagDir)

		manifest, err :=  getManifestForTag(tagDir)
		if err != nil {
			return err
		}
		r.t2m[tag] = manifest

		addTagToTagsList(r, manifest, tag)
	}

	return nil
}


func (r *Registry) GetManifestForTag(tag string) (string, error){
	err := r.buildTagAndManifestMaps()
	if err != nil {
		return "", fmt.Errorf("cannot build manifestsToTag map: %w", err)
	}

	manifest, ok := r.t2m[tag]
	if !ok {
		return "", fmt.Errorf("'%s' tag not found: manifest for tag does not exist", tag)
	}

	return manifest, nil
}

func (r *Registry) GetTagsWithManifest(manifest string) ([]string, error){
	err := r.buildTagAndManifestMaps()
	if err != nil {
		return nil, fmt.Errorf("cannot build tagsWithManifest map: %w", err)
	}

	tags, ok := r.m2t[manifest]
	if !ok {
		return nil, fmt.Errorf("something wrong: tagsWithManifest map does not contain manifest %s", manifest)
	}

	return tags, nil
}

func (r *Registry) GetTagsWithSameManifest(inTag string, excludeSelf bool) ([]string, error){
	err := r.buildTagAndManifestMaps()
	if err != nil {
		return nil, fmt.Errorf("cannot build tagsWithManifest map: %w", err)
	}

	manifest, err := r.GetManifestForTag(inTag)
	if err != nil {
		return nil, fmt.Errorf("cannot get manifest for tag: %w", err)
	}

	tags, ok := r.m2t[manifest]
	if !ok {
		return nil, fmt.Errorf("something wrong: tagsWithManifest map does not contain manifest %s", manifest)
	}


	if excludeSelf {
		tagsFiltered := make([]string, 0, len(tags) - 1)
		for _, tag := range tags {
			if tag != inTag {
				tagsFiltered = append(tagsFiltered, tag)
			}
		}
		tags = tagsFiltered
	}

	return tags, nil
}

func (r *Registry) IsTagSaveToRemove(tag string) (bool, error){
	manifest, err := r.GetManifestForTag(tag)
	if err != nil {
		return false, err
	}

	tags, err := r.GetTagsWithManifest(manifest)
	if err != nil {
		return false, err
	}

	return len(tags) == 1, nil
}