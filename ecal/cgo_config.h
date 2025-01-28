#ifndef ECAL_GO_CONFIG
#define ECAL_GO_CONFIG

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

void *NewConfig(void);
void DeleteConfig(void *config);

// Logging
void ConfigLoggingConsole(void* config, bool enable);
void ConfigLoggingConsoleAll(void* config);
void ConfigLoggingUdpReceive(void* config, bool enable);

void ConfigGetLoadedFilePath(uintptr_t handle);

// Publishers
bool ConfigPubShmEnabled(void);
bool ConfigPubUdpEnabled(void);
bool ConfigPubTcpEnabled(void);

// Subscribers
bool ConfigSubShmEnabled(void);
bool ConfigSubUdpEnabled(void);
bool ConfigSubTcpEnabled(void);


#ifdef __cplusplus
}
#endif

#endif
