package protobuf

import (
	"log"
	"slices"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// CloneOf returns a deep copy of m. If the top-level message is invalid,
// it returns an invalid message as well.
func CloneOf[M proto.Message](m M) M {
	return proto.Clone(m).(M)
}

func hasFile(fsetFile []*descriptorpb.FileDescriptorProto, fname string) bool {
	return slices.ContainsFunc(
		fsetFile,
		func(f *descriptorpb.FileDescriptorProto) bool { return f.GetName() == fname },
	)
}

func getFileDescriptor(desc protoreflect.MessageDescriptor, fset *descriptorpb.FileDescriptorSet) {
	if desc == nil {
		return
	}

	fdesc := desc.ParentFile()
	imports := fdesc.Imports()

	// TODO: Iterate services after the enums
	for i := range imports.Len() {
		sfdesc := imports.Get(i)

		// Iterate messages
		for m := range sfdesc.Messages().Len() {
			getFileDescriptor(sfdesc.Messages().Get(m), fset)
		}

		// Iterate enums
		if sfdesc.Enums().Len() > 0 {
			efdesc := sfdesc.Enums().Get(0).ParentFile()
			if !hasFile(fset.GetFile(), efdesc.Path()) {
				// Add the file to the set
				fset.File = append(fset.File, protodesc.ToFileDescriptorProto(efdesc))
			}
		}
	}

	if hasFile(fset.GetFile(), fdesc.Path()) {
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
