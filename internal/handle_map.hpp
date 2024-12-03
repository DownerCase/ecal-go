#ifndef ECAL_GO_HANDLE_MAP_HPP
#define ECAL_GO_HANDLE_MAP_HPP

#include <cstdint>
#include <unordered_map>

template <class T> class handle_map {
private:
  std::unordered_map<uintptr_t, T> handles{};

public:
  using iterator = typename decltype(handles)::iterator;
  template <typename... Args>
  std::pair<iterator, bool> emplace(uintptr_t handle, Args &&...args) {
    // Create a new T and store with the handle
    return this->handles.emplace(handle, T{std::forward<Args>(args)...});
  }

  bool erase(uintptr_t handle) { return this->handles.erase(handle) == 1; }

  T *find(uintptr_t handle) {
    const auto elem = this->handles.find(handle);
    if (elem == handles.end()) {
      return nullptr;
    }
    return &elem->second;
  }
};

#endif
