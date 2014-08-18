package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"models"
	"protodata"
)

func duplicateProtoList(duplicates []*models.DuplicateData, configs []*models.ConfigDuplicate) []*protodata.ChapterData {

	//list := models.ConfigDuplicateList()

	var result []*protodata.ChapterData
	result = append(result, &protodata.ChapterData{
		ChapterId:   proto.Int32(int32(configs[0].Chapter)),
		ChapterName: proto.String(configs[0].ChapterName),
		ChapterDesc: proto.String(""),
		IsUnlock:    proto.Bool(true),
	})

	for index, section := range configs {

		var sectionProto protodata.SectionData
		sectionProto.SectionId = proto.Int32(int32(section.Section))
		sectionProto.SectionName = proto.String(section.SectionName)
		sectionProto.SectionDesc = proto.String("")
		sectionProto.IsUnlock = proto.Bool(true)

		var find bool
		if index > 0 {
			for _, d := range duplicates {
				if d.Chapter == configs[index-1].Chapter && d.Section == configs[index-1].Section {
					find = true
					break
				} else {
					find = false
				}
			}
			if !find {
				sectionProto.IsUnlock = proto.Bool(false)
			}
		}

		if section.Chapter != int(*result[len(result)-1].ChapterId) {

			result = append(result, &protodata.ChapterData{
				ChapterId:   proto.Int32(int32(section.Chapter)),
				ChapterName: proto.String(section.ChapterName),
				ChapterDesc: proto.String(""),
				IsUnlock:    proto.Bool(find),
			})
		}
		result[len(result)-1].Sections = append(result[len(result)-1].Sections, &sectionProto)
	}

	return result
}
