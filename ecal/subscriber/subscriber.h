#ifndef ECAL_GO_SUBSCRIBER_H
#define ECAL_GO_SUBSCRIBER_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

const void *NewSubscriber();
bool DestroySubscriber(uintptr_t handle);

bool SubscriberCreate(
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

// Receive a message and return a handle to it as well as the message data
// pointer and length
uintptr_t
SubscriberReceive(uintptr_t subscriber_handle, const char **msg, size_t *len);

bool DestroyMessage(uintptr_t message_handle);

#ifdef __cplusplus
}
#endif

#endif
