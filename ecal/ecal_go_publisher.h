#ifndef ECAL_GO_PUBLISHER_H
#define ECAL_GO_PUBLISHER_H

#include <stdbool.h>
#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

const void *NewPublisher();
bool DestroyPublisher(uintptr_t handle);

bool PublisherCreate(uintptr_t handle, const char* const topic);
void PublisherSend(uintptr_t handle, void* buf, size_t len);

#ifdef __cplusplus
}
#endif

#endif
