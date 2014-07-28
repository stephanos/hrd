package hrd

import (
	_ "github.com/101loops/bdd"
)

//var _ = Describe("Codec", func() {
//
//	It("return simple codec", func() {
//		validateSimpleCode()
//		validateSimpleCode() // ... coming from cache
//	})
//
//	It("return complex codec", func() {
//		code, err := getCodec(&ComplexModel{})
//
//		Check(err, IsNil)
//		Check(code, NotNil)
//		Check(code.complete, IsTrue)
//
//		Check(code.byIndex, Equals, map[int]tagCodec{
//			0: tagCodec{name: "tag", tags: []string{}},
//			//1: tagCodec{name: "pair", tags: []string{}},
//			1: tagCodec{name: "tags", tags: []string{}},
//			//3: tagCodec{name: "pairs", tags: []string{}},
//		})
//
//		Check(code.byName, HasLen, 4)
//		Check(code.byName, HasKeys, "tag.key", "tag.val", "tags.key", "tags.val")
//		//"pair.key", "pair.val", "pairs.key", "pairs.val"
//
//		Check(code.byName["tag.key"].subcodec, NotNil)
//		Check(*code.byName["tag.key"].subcodec, Equals, pairCodec)
//		Check(code.byName["tag.val"].subcodec, NotNil)
//		Check(*code.byName["tag.val"].subcodec, Equals, pairCodec)
//
//		//Check(code.byName["pair.key"].subcodec, IsNil)
//		//Check(*code.byName["pair.key"].subcodec, Equals,  pairCodec)
//		//Check(code.byName["pair.val"].subcodec, IsNil)
//		//Check(*code.byName["pair.val"].subcodec, Equals,  pairCodec)
//
//		Check(code.byName["tags.key"].subcodec, NotNil)
//		Check(*code.byName["tags.key"].subcodec, Equals, pairCodec)
//		Check(code.byName["tags.val"].subcodec, NotNil)
//		Check(*code.byName["tags.val"].subcodec, Equals, pairCodec)
//
//		//Check(code.byName["pairs.key"].subcodec, IsNil)
//		//Check(*code.byName["pairs.key"].subcodec, Equals,  pairCodec)
//		//Check(code.byName["pairs.val"].subcodec, IsNil)
//		//Check(*code.byName["pairs.val"].subcodec, Equals,  pairCodec)
//
//		Check(code.hasSlice, IsTrue)
//	})
//})
//
//func validateSimpleCode() {
//	code, err := getCodec(&SimpleModel{})
//
//	Check(err, IsNil)
//	Check(code, NotNil)
//	Check(*code, Equals, hrdCodec{
//		byIndex: map[int]*{
//			1: tagCodec{name: "num", tags: []string{}},
//			3: tagCodec{name: "dat", tags: []string{"index"}},
//			4: tagCodec{name: "html", tags: []string{"index", "omitempty"}},
//			5: tagCodec{name: "timing", tags: []string{"index", "omitempty"}},
//		},
//		byName: map[string]fieldCodec{
//			"num":    fieldCodec{index: 1, subcodec: nil},
//			"dat":    fieldCodec{index: 3, subcodec: nil},
//			"html":   fieldCodec{index: 4, subcodec: nil},
//			"timing": fieldCodec{index: 5, subcodec: nil},
//		},
//		hasSlice: false,
//		complete: true,
//	})
//}
//
//var (
//	pairCodec = hrdCodec{
//		byIndex: map[int]tagCodec{
//			0: tagCodec{name: "key", tags: []string{"index", "omitempty"}},
//			1: tagCodec{name: "val", tags: []string{}},
//		},
//		byName: map[string]fieldCodec{
//			"key": fieldCodec{index: 0, subcodec: nil},
//			"val": fieldCodec{index: 1, subcodec: nil},
//		},
//		hasSlice: false,
//		complete: true,
//	}
//)
