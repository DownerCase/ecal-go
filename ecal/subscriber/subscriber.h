#ifndef ECAL_GO_SUBSCRIBER_H
#define ECAL_GO_SUBSCRIBER_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

bool NewSubscriber(uintptr_t handle);
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

#ifdef __cplusplus
}
#endif

#endif