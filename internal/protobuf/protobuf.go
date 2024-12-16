package protobuf

import (
	"log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func hasFile(fset *descriptorpb.FileDescriptorSet, fname string) bool {
	for _, file := range fset.GetFile() {
		if file.GetName() == fname {
			return true
		}
	}
	return false
}

func getFileDescriptor(desc protoreflect.MessageDescriptor, fset *descriptorpb.FileDescriptorSet) {
	if desc == nil {
		return
	}

	fdesc := desc.ParentFile()
	imports := fdesc.Imports()

	for i := range imports.Len() {
		sfdesc := imports.Get(i)

		// Iterate messages
		for m := range sfdesc.Messages().Len() {
			getFileDescriptor(sfdesc.Messages().Get(m), fset)
		}

		// Iterate enums
		if sfdesc.Enums().Len() > 0 {
			edesc := sfdesc.Enums().Get(0)
			efdesc := edesc.ParentFile()
			if !hasFile(fset, efdesc.Path()) {
				// Add the file to the set
				fset.File = append(fset.File, protodesc.ToFileDescriptorProto(efdesc))
			}
		}

		// TODO: Iterate services
	}

	if hasFile(fset, fdesc.Path()) {
		// File already added to descriptor set, continue
		return
	}
	fset.File = append(fset.File, protodesc.ToFileDescriptorProto(fdesc))

	// Add fields
	for f := range desc.Fields().Len() {
		getFileDescriptor(desc.Fields().Get(f).Message(), fset)
	}
}

func GetProtoMessageDescription(msg proto.Message) []byte {
	desc := msg.ProtoReflect().Descriptor()
	pset := descriptorpb.FileDescriptorSet{}
	getFileDescriptor(desc, &pset)
	bytes, err := proto.Marshal(&pset)
	if err != nil {
		log.Println("WARN: GetProtoMessageDescription failed to marshal file descriptor set", err)
		return nil
	}
	return bytes
}

func GetFullName(msg proto.Message) string {
	return string(msg.ProtoReflect().Descriptor().FullName())
}
