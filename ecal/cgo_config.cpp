#include "cgo_config.h"

#include <ecal/config.h>
#include <ecal/init.h>

extern "C" {

extern void goCopyString(uintptr_t, const char *);

void *NewConfig(void) {
  auto *config = new eCAL::Configuration;
  *config      = eCAL::Init::Configuration();
  return config;
}

void DeleteConfig(void *config) {
  delete reinterpret_cast<eCAL::Configuration *>(config);
}

void ConfigLoggingConsole(void *config, bool enable) {
  auto *cfg = reinterpret_cast<eCAL::Configuration *>(config);
  cfg->logging.provider.console.enable = enable;
}

void ConfigLoggingConsoleAll(void *config) {
  auto *cfg = reinterpret_cast<eCAL::Configuration *>(config);
  cfg->logging.provider.console.filter_log =
      eCAL::Logging::eLogLevel::log_level_all;
}

void ConfigLoggingUdpReceive(void *config, bool enable) {
  auto *cfg = reinterpret_cast<eCAL::Configuration *>(config);
  cfg->logging.receiver.enable = enable;
}

void ConfigGetLoadedFilePath(uintptr_t handle) {
  const auto cfg = eCAL::Config::GetLoadedEcalIniPath();
  goCopyString(handle, cfg.c_str());
}

bool ConfigPubShmEnabled(void) {
  return eCAL::GetPublisherConfiguration().layer.shm.enable;
}
bool ConfigPubUdpEnabled(void) {
  return eCAL::GetPublisherConfiguration().layer.udp.enable;
}
bool ConfigPubTcpEnabled(void) {
  return eCAL::GetPublisherConfiguration().layer.tcp.enable;
}

bool ConfigSubShmEnabled(void) {
  return eCAL::GetSubscriberConfiguration().layer.shm.enable;
}
bool ConfigSubUdpEnabled(void) {
  return eCAL::GetSubscriberConfiguration().layer.udp.enable;
}
bool ConfigSubTcpEnabled(void) {
  return eCAL::GetSubscriberConfiguration().layer.tcp.enable;
}
}
