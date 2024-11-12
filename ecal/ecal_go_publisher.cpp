#include "ecal_go_publisher.h"

#include <memory>
#include <unordered_map>

#include <ecal/ecal_publisher.h>

namespace {
std::unordered_map<uintptr_t, std::unique_ptr<eCAL::CPublisher>> publishers;

eCAL::CPublisher *const getPublisher(uintptr_t handle) {
  const auto publisher = publishers.find(handle);
  if (publisher == publishers.end()) {
    return nullptr;
  }
  return publisher->second.get();
}

} // namespace

const void *NewPublisher() {
  auto publisher = std::make_unique<eCAL::CPublisher>();
  const auto handle = publisher.get();
  const auto [new_pub, added] = publishers.emplace(
      reinterpret_cast<uintptr_t>(handle), std::move(publisher));
  if (!added) {
    return nullptr;
  }
  return handle;
}

bool DestroyPublisher(uintptr_t handle) {
  return publishers.erase(handle) == 1;
}

bool PublisherCreate(uintptr_t handle, const char *const topic) {
  auto *publisher = getPublisher(handle);
  if (publisher == nullptr) {
    return false;
  }
  return publisher->Create(topic, {"std::string", "base"});
}

void PublisherSend(uintptr_t handle, void *buf, size_t len) {
  auto *publisher = getPublisher(handle);
  if (publisher == nullptr) {
    return;
  }
  publisher->Send(buf, len);
}
