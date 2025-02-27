// This example demonstrates how to set and store metadata using GStreamer.
//
// Some elements support setting tags on a media stream. An example would be
// id3v2mux. The element signals this by implementing The GstTagsetter interface.
// You can query any element implementing this interface from the pipeline, and
// then tell the returned implementation of GstTagsetter what tags to apply to
// the media stream.
//
// This example's pipeline creates a new flac file from the testaudiosrc
// that the example application will add tags to using GstTagsetter.
// The operated pipeline looks like this:
//
//   {audiotestsrc} - {flacenc} - {filesink}
//
// For example for pipelines that transcode a multimedia file, the input
// already has tags. For cases like this, the GstTagsetter has the merge
// setting, which the application can configure to tell the element
// implementing the interface whether to merge newly applied tags to the
// already existing ones, or if all existing ones should replace, etc.
// (More modes of operation are possible, see: gst.TagMergeMode)
// This merge-mode can also be supplied to any method that adds new tags.
package main

import (
	"fmt"
	"time"

	"github.com/clintlombard/go-gst/examples"
	"github.com/clintlombard/go-gst/gst"
)

func tagsetter() error {
	gst.Init(nil)

	pipeline, err := gst.NewPipelineFromString(
		"audiotestsrc wave=white-noise num-buffers=10000 ! flacenc ! filesink location=test.flac",
	)
	if err != nil {
		return err
	}

	// Query the pipeline for elements implementing the GstTagsetter interface.
	// In our case, this will return the flacenc element.
	element, err := pipeline.GetByInterface(gst.InterfaceTagSetter)
	if err != nil {
		return err
	}

	// We actually just retrieved a *gst.Element with the above call. We can retrieve
	// the underying TagSetter interface like this.
	tagsetter := element.TagSetter()

	// Tell the element implementing the GstTagsetter interface how to handle already existing
	// metadata.
	tagsetter.SetTagMergeMode(gst.TagMergeKeepAll)

	// Set the "title" tag to "Special randomized white-noise".
	//
	// The first parameter gst.TagMergeAppend tells the tagsetter to append this title
	// if there already is one.
	tagsetter.AddTagValue(gst.TagMergeAppend, gst.TagTitle, "Special randomized white-noise")

	pipeline.SetState(gst.StatePlaying)

	var cont bool
	var pipelineErr error
	for {
		msg := pipeline.GetPipelineBus().TimedPop(time.Duration(-1))
		if msg == nil {
			break
		}
		if cont, pipelineErr = handleMessage(msg); pipelineErr != nil || !cont {
			pipeline.SetState(gst.StateNull)
			break
		}
	}

	return pipelineErr
}

func handleMessage(msg *gst.Message) (bool, error) {
	defer msg.Unref()
	switch msg.Type() {
	case gst.MessageTag:
		fmt.Println(msg) // Prirnt our tags
	case gst.MessageEOS:
		return false, nil
	case gst.MessageError:
		return false, msg.ParseError()
	}
	return true, nil
}

func main() {
	examples.Run(tagsetter)
}
