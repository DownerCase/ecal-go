#ifndef ECAL_GO_REGISTRATION
#define ECAL_GO_REGISTRATION

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

size_t AddPublisherEventCallback(uintptr_t handle);
void RemPublisherEventCallback(size_t handle);

size_t AddSubscriberEventCallback(uintptr_t handle);
void RemSubscriberEventCallback(size_t handle);

#ifdef __cplusplus
}
#endif

#endif
