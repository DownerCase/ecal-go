// Implementation for cgo preamble functions
package publisher

// #include "publisher.h"
// void *GoNewPublisher(
//  _GoString_ topic,
//  _GoString_ name, _GoString_ encoding,
//  const char* const descriptor, size_t descriptor_len
// ) {
//  return NewPublisher(
//    _GoStringPtr(topic), _GoStringLen(topic),
//    _GoStringPtr(name), _GoStringLen(name),
//    _GoStringPtr(encoding), _GoStringLen(encoding),
//    descriptor, descriptor_len
//  );
//}
//// C preamble.
import "C"
