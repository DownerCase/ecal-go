#ifndef ECAL_GO_REGISTRATION
#define ECAL_GO_REGISTRATION

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct CEntityId {
  uint64_t entity_id;
  const char *host_name;
  int32_t process_id;
  char PADDING[4];
};

struct CTopicId {
  struct CEntityId topic_id;
  const char *topic_name;
};

size_t AddPublisherEventCallback(uintptr_t handle);
void RemPublisherEventCallback(size_t handle);

size_t AddSubscriberEventCallback(uintptr_t handle);
void RemSubscriberEventCallback(size_t handle);

#ifdef __cplusplus
}
#endif

#endif
