package models_test

import (
	"fmt"

	"cuboid-challenge/app/models"
	"cuboid-challenge/app/tests/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bag", func() {
	testutils.LoadEnv()

	Describe("Validations", func() {
		var bag models.Bag

		BeforeEach(func() {
			bag = models.Bag{
				Title:   "My Bag",
				Volume:  10,
				Cuboids: nil,
			}
		})

		It("Requires a Title", func() {
			bag.Title = ""
			isValid, validationErr := models.Validate(bag)

			Expect(isValid).To(Equal(false))
			Expect(validationErr[0].Field).To(Equal("Title"))
			Expect(validationErr[0].Type).To(Equal("required"))
		})

		It("Requires a Volume", func() {
			bag.Volume = 0
			isValid, validationErr := models.Validate(bag)

			Expect(isValid).To(Equal(false))
			Expect(validationErr[0].Field).To(Equal("Volume"))
			Expect(validationErr[0].Type).To(Equal("gt"))
			Expect(validationErr[0].Param).To(Equal("0"))
		})
	})

	bags := map[string]struct {
		Bag                  models.Bag
		expectedPayloadVol   uint
		expectedAvailableVol uint
	}{
		"no cuboids": {
			Bag: models.Bag{
				Volume: 10,
				Title:  "A bag with no cuboids",
			},
			expectedPayloadVol:   0,
			expectedAvailableVol: 10,
		},
		"one cuboids": {
			Bag: models.Bag{
				Volume: 20,
				Title:  "A bag with one cuboid",
				Cuboids: []models.Cuboid{
					{Width: 3, Height: 2, Depth: 3},
				},
			},
			expectedPayloadVol:   18,
			expectedAvailableVol: 2,
		},
	}
	for key, v := range bags {
		value := v
		Describe(fmt.Sprintf("with %s", key), func() {
			It(fmt.Sprintf("Has %v PayloadVolume", value.expectedPayloadVol), func() {
				Expect(value.Bag.PayloadVolume()).To(Equal(value.expectedPayloadVol))
			})

			It(fmt.Sprintf("Has %v AvailableVolume", value.expectedAvailableVol), func() {
				Expect(value.Bag.AvailableVolume()).To(BeEquivalentTo(value.expectedAvailableVol))
			})

			Specify("To JSON", func() {
				m, err := testutils.Serialize(&value.Bag)
				Expect(err).NotTo(HaveOccurred())
				Expect(m["volume"]).To(BeEquivalentTo(value.Bag.Volume))
				Expect(m["title"]).To(BeEquivalentTo(value.Bag.Title))
				Expect(m["payloadVolume"]).To(BeEquivalentTo(value.expectedPayloadVol))
				Expect(m["availableVolume"]).To(BeEquivalentTo(value.expectedAvailableVol))
			})
		})
	}
})
