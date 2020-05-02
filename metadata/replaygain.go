// Copyright (c) 2018-2019 ChrisOboe <chris@oboe.email>
// SPDX-License-Identifier: GPL-3.0

package metadata

/*
import (
	"github.com/giorgisio/goav/avfilter"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
	"github.com/pkg/errors"
)

type Ebur128 struct {
	Integrated        float64
	Peak              float64
	LRA               float64
	Threshold         float64
	NormalizationType string
	TargetOffset      float64
}

func CalcEbur128(file string) (Ebur128, error) {
	// Open file
	pFormatContext := avformat.AvformatAllocContext()
	errCode := avformat.AvformatOpenInput(&pFormatContext, file, nil, nil)
	if errCode != 0 {
		return Ebur128{}, errors.New("Couldn't open file")
	}
	defer pFormatContext.AvformatCloseInput()

	errCode = pFormatContext.AvformatFindStreamInfo(nil)
	if errCode < 0 {
		return Ebur128{}, errors.New("Couldn't get stream info")
	}

	// create the graph
	// (input) -> abuffer -> loudness -> (output /dev/null)
	pGraph := avfilter.AvfilterGraphAlloc()
	defer pGraph.AvfilterGraphFree()

	// create the abuffer filter
	abuffer := avfilter.AvfilterGetByName("abuffer")
	abuffer.AvfilterRegister()
	abuffer_ctx := pGraph.AvfilterGraphAllocFilter(abuffer, "src")

	// create the loudness filter
	loudness := avfilter.AvfilterGetByName("loudness")
	loudness.AvfilterRegister()

	return Ebur128{}
}

*/
