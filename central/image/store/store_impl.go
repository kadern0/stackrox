package store

import (
	"fmt"
	"time"

	"bitbucket.org/stack-rox/apollo/central/globaldb/ops"
	"bitbucket.org/stack-rox/apollo/central/metrics"
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	imagesPkg "bitbucket.org/stack-rox/apollo/pkg/images"
	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
)

type storeImpl struct {
	db *bolt.DB
}

// ListImage returns ListImage with given sha.
func (b *storeImpl) ListImage(sha string) (image *v1.ListImage, exists bool, err error) {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.Get, "ListImage")

	digest := imagesPkg.NewDigest(sha).Digest()
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(listImageBucket))
		image = new(v1.ListImage)
		val := bucket.Get([]byte(digest))
		if val == nil {
			return nil
		}
		exists = true
		return proto.Unmarshal(val, image)
	})
	return
}

// ListImages returns all ListImages
func (b *storeImpl) ListImages() (images []*v1.ListImage, err error) {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.GetMany, "ListImage")

	err = b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(listImageBucket))
		return b.ForEach(func(k, v []byte) error {
			var image v1.ListImage
			if err := proto.Unmarshal(v, &image); err != nil {
				return err
			}
			images = append(images, &image)
			return nil
		})
	})
	return
}

// GetImages returns all images regardless of request
func (b *storeImpl) GetImages() (images []*v1.Image, err error) {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.GetAll, "Image")

	err = b.db.View(func(tx *bolt.Tx) error {
		images, err = readAllImages(tx)
		return err
	})
	return
}

// CountImages returns all images regardless of request
func (b *storeImpl) CountImages() (count int, err error) {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.GetAll, "Image")

	err = b.db.View(func(tx *bolt.Tx) error {
		count, err = countAllImages(tx)
		return err
	})
	return
}

// GetImage returns image with given sha.
func (b *storeImpl) GetImage(sha string) (image *v1.Image, exists bool, err error) {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.Get, "Image")

	err = b.db.View(func(tx *bolt.Tx) error {
		exists = hasImage(tx, []byte(idForSha(sha)))
		if !exists {
			return nil
		}
		image, err = readImage(tx, []byte(idForSha(sha)))
		return err
	})
	return
}

// GetImagesBatch returns image with given sha.
func (b *storeImpl) GetImagesBatch(shas []string) (images []*v1.Image, err error) {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.GetMany, "Image")

	err = b.db.View(func(tx *bolt.Tx) error {
		for _, sha := range shas {
			image, err := readImage(tx, []byte(idForSha(sha)))
			if err != nil {
				return err
			}
			images = append(images, image)
		}
		return nil
	})
	return
}

// UpdateImage updates a image to bolt.
func (b *storeImpl) UpsertImage(image *v1.Image) error {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.Upsert, "Image")

	return b.db.Update(func(tx *bolt.Tx) error {
		err := writeImage(tx, image)
		if err != nil {
			return err
		}
		return upsertListImage(tx, image)
	})
}

// DeleteImage deletes an image an all it's data (but maintains the orch sha to registry sha mapping).
func (b *storeImpl) DeleteImage(sha string) error {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.Remove, "Image")

	return b.db.Update(func(tx *bolt.Tx) error {
		err := deleteImage(tx, []byte(idForSha(sha)))
		if err != nil {
			return err
		}
		return deleteListImage(tx, []byte(idForSha(sha)))
	})
}

// GetRegistrySha retrieves a sha to registry sha mapping.
func (b *storeImpl) GetRegistrySha(orchSha string) (regSha string, exists bool, err error) {
	defer metrics.SetBoltOperationDurationTime(time.Now(), ops.Get, "ImageRegistrySha")

	err = b.db.View(func(tx *bolt.Tx) error {
		exists = hasSha(tx, []byte(idForSha(orchSha)))
		if !exists {
			return nil
		}
		regSha, err = readSha(tx, []byte(idForSha(orchSha)))
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// UpsertRegistrySha adds a sha to registry sha mapping.
func (b *storeImpl) UpsertRegistrySha(orchSha string, regSha string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return writeSha(tx, []byte(idForSha(orchSha)), regSha)
	})
}

// DeleteRegistrySha removes a sha to registry sha mapping.
func (b *storeImpl) DeleteRegistrySha(orchSha string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return deleteSha(tx, []byte(idForSha(orchSha)))
	})
}

// General helper functions.
////////////////////////////

func idForSha(sha string) string {
	return imagesPkg.NewDigest(sha).Digest()
}

func convertImageToListImage(i *v1.Image) *v1.ListImage {
	listImage := &v1.ListImage{
		Sha:     i.GetName().GetSha(),
		Name:    i.GetName().GetFullName(),
		Created: i.GetMetadata().GetCreated(),
	}

	if i.GetScan() != nil {
		listImage.SetComponents = &v1.ListImage_Components{
			Components: int64(len(i.GetScan().GetComponents())),
		}
		var numVulns int64
		var numFixableVulns int64
		for _, c := range i.GetScan().GetComponents() {
			numVulns += int64(len(c.GetVulns()))
			for _, v := range c.GetVulns() {
				if v.FixedBy != "" {
					numFixableVulns++
				}
			}
		}
		listImage.SetCves = &v1.ListImage_Cves{
			Cves: numVulns,
		}
		listImage.SetFixable = &v1.ListImage_FixableCves{
			FixableCves: numFixableVulns,
		}
	}
	return listImage
}

// In-Transaction helper functions.
///////////////////////////////////

// readAllImages reads all the images in the DB within a transaction.
func readAllImages(tx *bolt.Tx) (images []*v1.Image, err error) {
	bucket := tx.Bucket([]byte(imageBucket))
	err = bucket.ForEach(func(k, v []byte) error {
		image, err := readImage(tx, k)
		if err != nil {
			return err
		}

		images = append(images, image)
		return nil
	})
	return
}

// readAllImages reads all the images in the DB within a transaction.
func countAllImages(tx *bolt.Tx) (count int, err error) {
	bucket := tx.Bucket([]byte(imageBucket))
	count = 0
	err = bucket.ForEach(func(k, v []byte) error {
		count++
		return nil
	})
	return
}

// HasImage returns whether a image exists for the given id.
func hasImage(tx *bolt.Tx, id []byte) bool {
	bucket := tx.Bucket([]byte(imageBucket))

	bytes := bucket.Get(id)
	if bytes == nil {
		return false
	}
	return true
}

// HasImage returns whether a image exists for the given id.
func hasSha(tx *bolt.Tx, id []byte) bool {
	bucket := tx.Bucket([]byte(orchShaToRegShaBucket))

	bytes := bucket.Get(id)
	if bytes == nil {
		return false
	}
	return true
}

// readImage reads a image within a transaction.
func readImage(tx *bolt.Tx, id []byte) (image *v1.Image, err error) {
	bucket := tx.Bucket([]byte(imageBucket))

	bytes := bucket.Get(id)
	if bytes == nil {
		err = fmt.Errorf("image with id: %s does not exist", id)
		return
	}

	image = new(v1.Image)
	err = proto.Unmarshal(bytes, image)
	return
}

// writeImage writes an image within a transaction.
func writeImage(tx *bolt.Tx, image *v1.Image) (err error) {
	bucket := tx.Bucket([]byte(imageBucket))

	id := []byte(idForSha(image.GetName().GetSha()))

	bytes, err := proto.Marshal(image)
	if err != nil {
		return
	}
	bucket.Put(id, bytes)
	return
}

// deleteImage deletes an image within a transaction.
func deleteImage(tx *bolt.Tx, id []byte) (err error) {
	bucket := tx.Bucket([]byte(imageBucket))

	bucket.Delete(id)
	return
}

// readSha reads a image within a transaction.
func readSha(tx *bolt.Tx, id []byte) (regSha string, err error) {
	bucket := tx.Bucket([]byte(orchShaToRegShaBucket))

	bytes := bucket.Get(id)
	if bytes == nil {
		err = fmt.Errorf("image with id: %s does not exist", id)
		return
	}

	regSha = string(bytes)
	return
}

// writeSha writes an image's sha within a transaction.
func writeSha(tx *bolt.Tx, id []byte, regSha string) (err error) {
	bucket := tx.Bucket([]byte(orchShaToRegShaBucket))

	bucket.Put(id, []byte(regSha))
	return
}

// deleteSha deletes an image's sha within a transaction.
func deleteSha(tx *bolt.Tx, id []byte) (err error) {
	bucket := tx.Bucket([]byte(orchShaToRegShaBucket))

	bucket.Delete(id)
	return
}

func upsertListImage(tx *bolt.Tx, image *v1.Image) error {
	bucket := tx.Bucket([]byte(listImageBucket))
	listImage := convertImageToListImage(image)
	bytes, err := proto.Marshal(listImage)
	if err != nil {
		return err
	}
	digest := imagesPkg.NewDigest(image.GetName().GetSha()).Digest()
	return bucket.Put([]byte(digest), bytes)
}

func deleteListImage(tx *bolt.Tx, id []byte) (err error) {
	bucket := tx.Bucket([]byte(listImageBucket))

	bucket.Delete(id)
	return
}
