#ifndef ECAL_GO_HANDLE_MAP_HPP
#define ECAL_GO_HANDLE_MAP_HPP

#include <cassert>
#include <cstdint>
#include <memory>
#include <type_traits>
#include <unordered_map>
#include <utility>

template <class T> class handle_map {
private:
  std::unordered_map<uintptr_t, std::unique_ptr<T>> handles;
public:
  using iterator = typename decltype(handles)::iterator;
  template <typename... Args>
  std::pair<iterator, bool> emplace(Args &&...args) {
    // Create a new T to get it's address to use as the key
    auto instance = std::make_unique<T>(std::forward<Args>(args)...);
    const auto &handle = *instance;
    static_assert(
        std::is_same_v<decltype(handle), const T &>,
        "Unable to take reference"
    );
    // Store in the map
    return this->handles.emplace(
        reinterpret_cast<std::uintptr_t>(&handle),
        std::move(instance)
    );
  }

  bool erase(std::uintptr_t handle) { return this->handles.erase(handle) == 1; }

  T *find(std::uintptr_t handle) {
    const auto elem = this->handles.find(handle);
    if (elem == handles.end()) {
      return nullptr;
    }
    return elem->second.get();
  }

};

#endif
