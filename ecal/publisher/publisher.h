#ifndef ECAL_GO_PUBLISHER_H
#define ECAL_GO_PUBLISHER_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

bool NewPublisher(uintptr_t handle);
bool DestroyPublisher(uintptr_t handle);

bool PublisherCreate(
    uintptr_t handle,
    const char *const topic,
    size_t topic_len,
    const char *const datatype_name,
    size_t datatype_name_len,
    const char *const datatype_encoding,
    size_t datatype_encoding_len,
    const char *const datatype_descriptor,
    size_t datatype_descriptor_len
);
void PublisherSend(uintptr_t handle, void *buf, size_t len);

#ifdef __cplusplus
}
#endif

#endif
