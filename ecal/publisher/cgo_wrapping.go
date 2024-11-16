// Implementation for cgo preamble functions
package publisher

// #include "publisher.h"
//bool GoPublisherCreate(
//  uintptr_t handle,
//  _GoString_ topic,
//  _GoString_ name, _GoString_ encoding,
//  const char* const descriptor, size_t descriptor_len
//) {
//  return PublisherCreate(
//    handle,
//    _GoStringPtr(topic), _GoStringLen(topic),
//    _GoStringPtr(name), _GoStringLen(name),
//    _GoStringPtr(encoding), _GoStringLen(encoding),
//    descriptor, descriptor_len
//  );
//}
import "C"
