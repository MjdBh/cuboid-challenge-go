package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"cuboid-challenge/app/tests/factories"
	"cuboid-challenge/app/tests/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bag Controller", func() {
	testutils.LoadEnv()
	testutils.ConnectDB()
	testutils.ClearDB()

	AfterEach(func() {
		testutils.ClearDB()
	})

	var w *httptest.ResponseRecorder

	Describe("List", func() {
		BeforeEach(func() {
			for i := 0; i < 3; i++ {
				testutils.AddRecords(factories.Bag())
			}
			w = testutils.MockRequest(http.MethodGet, "/bags", nil)
		})

		It("Response HTTP status code 200", func() {
			Expect(w.Code).To(Equal(200))
		})

		It("Fetch all bags", func() {
			l, _ := testutils.DeserializeList(w.Body.String())
			Expect(len(l)).To(Equal(3))
			for _, m := range l {
				Expect(m["title"]).ToNot(BeNil())
				Expect(m["volume"]).ToNot(BeNil())
				Expect(m["cuboids"]).ToNot(BeNil())
			}
		})
	})

	Describe("Get", func() {
		var bagID uint
		bag := factories.Bag()

		BeforeEach(func() {
			testutils.AddRecords(&bag)
		})

		JustBeforeEach(func() {
			w = testutils.MockRequest(http.MethodGet, "/bags/"+fmt.Sprintf("%v", bagID), nil)
		})

		Context("When the bag is present", func() {
			BeforeEach(func() {
				bagID = bag.ID
			})

			It("Response HTTP status code 200", func() {
				Expect(w.Code).To(Equal(200))
			})

			It("Get the bag", func() {
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["id"]).To(BeEquivalentTo(bag.ID))
				Expect(m["title"]).ToNot(BeNil())
				Expect(m["volume"]).ToNot(BeNil())
				Expect(m["cuboids"]).ToNot(BeNil())
			})
		})

		Context("When the bag is not present", func() {
			BeforeEach(func() {
				bagID = 99999
			})

			It("Response HTTP status code 404", func() {
				Expect(w.Code).To(Equal(404))
			})
		})
	})

	Describe("Create", func() {
		bagPayload := map[string]interface{}{}

		JustBeforeEach(func() {
			body, _ := testutils.SerializeToString(bagPayload)
			w = testutils.MockRequest(http.MethodPost, "/bags", &body)
		})

		Context("When the payload is valid", func() {
			BeforeEach(func() {
				bagPayload["title"] = "My title"
				bagPayload["volume"] = 10
			})

			It("Response HTTP status code 201", func() {
				Expect(w.Code).To(Equal(201))
			})

			It("Create the bag", func() {
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["id"]).ToNot(BeNil())
				Expect(m["title"]).ToNot(BeNil())
				Expect(m["volume"]).ToNot(BeNil())
				Expect(m["cuboids"]).To(BeEmpty())
			})
		})

		Context("When the payload is invalid", func() {
			BeforeEach(func() {
				bagPayload["title"] = ""
				bagPayload["volume"] = 0
			})

			It("Response HTTP status code 400", func() {
				Expect(w.Code).To(Equal(400))
			})
		})
	})
})
