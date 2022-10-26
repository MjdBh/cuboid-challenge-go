package models_test

import (
	"fmt"

	"cuboid-challenge/app/models"
	"cuboid-challenge/app/tests/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cuboid", func() {
	testutils.LoadEnv()

	Describe("Validations", func() {
		var cuboid models.Cuboid

		BeforeEach(func() {
			cuboid = models.Cuboid{
				Width:  2,
				Height: 3,
				Depth:  4,
				BagID:  10,
			}
		})

		validateGtThan0 := func(prop string) {
			isValid, validationErr := models.Validate(cuboid)

			Expect(isValid).To(Equal(false))
			Expect(validationErr[0].Field).To(Equal(prop))
			Expect(validationErr[0].Type).To(Equal("gt"))
			Expect(validationErr[0].Param).To(Equal("0"))
		}

		It("Requires a Width", func() {
			cuboid.Width = 0
			validateGtThan0("Width")
		})

		It("Requires a Height", func() {
			cuboid.Height = 0
			validateGtThan0("Height")
		})

		It("Requires a Depth", func() {
			cuboid.Depth = 0
			validateGtThan0("Depth")
		})
	})

	cuboids := map[string]struct {
		Cuboid      models.Cuboid
		expectedVol uint
	}{
		"3 x 3 x 3": {
			Cuboid:      models.Cuboid{Width: 3, Height: 3, Depth: 3},
			expectedVol: 27,
		},
	}
	for key, v := range cuboids {
		value := v
		Describe(fmt.Sprintf("with %s", key), func() {
			It(fmt.Sprintf("Has %v volume", value.expectedVol), func() {
				Expect(value.Cuboid.PayloadVolume()).To(Equal(value.expectedVol))
			})

			Specify("To JSON", func() {
				m, err := testutils.Serialize(&value.Cuboid)
				Expect(err).NotTo(HaveOccurred())
				Expect(m["width"]).To(BeEquivalentTo(value.Cuboid.Width))
				Expect(m["height"]).To(BeEquivalentTo(value.Cuboid.Height))
				Expect(m["depth"]).To(BeEquivalentTo(value.Cuboid.Depth))
				Expect(m["volume"]).To(BeEquivalentTo(value.expectedVol))
				Expect(m["bagId"]).To(BeEquivalentTo(value.Cuboid.BagID))
			})
		})
	}
})
