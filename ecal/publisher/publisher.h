#ifndef ECAL_GO_PUBLISHER_H
#define ECAL_GO_PUBLISHER_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

void *NewPublisher(
    const char *topic,
    size_t topic_len,
    const char *datatype_name,
    size_t datatype_name_len,
    const char *datatype_encoding,
    size_t datatype_encoding_len,
    const char *datatype_descriptor,
    size_t datatype_descriptor_len
);
void DestroyPublisher(void* publisher);

void PublisherSend(void* publisher, void *buf, size_t len);

#ifdef __cplusplus
}
#endif

#endif
