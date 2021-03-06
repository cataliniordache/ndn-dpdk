#ifndef NDNDPDK_PINGSERVER_SERVER_H
#define NDNDPDK_PINGSERVER_SERVER_H

/** @file */

#include "../dpdk/thread.h"
#include "../iface/face.h"
#include "../iface/pktqueue.h"
#include "../vendor/pcg_basic.h"

#define PINGSERVER_MAX_PATTERNS 256
#define PINGSERVER_MAX_REPLIES 8
#define PINGSERVER_MAX_SUM_WEIGHT 256

typedef uint8_t PingReplyId;

typedef enum PingServerReplyKind
{
  PINGSERVER_REPLY_DATA,
  PINGSERVER_REPLY_NACK,
  PINGSERVER_REPLY_TIMEOUT,
} PingServerReplyKind;

typedef struct PingServerReply
{
  uint64_t nInterests;
  DataGen* dataGen;
  uint8_t kind;
  uint8_t nackReason;
} PingServerReply;

/** @brief Per-prefix information in ndnping server. */
typedef struct PingServerPattern
{
  LName prefix;
  uint16_t nReplies;
  uint16_t nWeights;
  PingReplyId weight[PINGSERVER_MAX_SUM_WEIGHT];
  PingServerReply reply[PINGSERVER_MAX_REPLIES];
  uint8_t prefixBuffer[NameMaxLength];
} PingServerPattern;

/** @brief ndnping server. */
typedef struct PingServer
{
  PktQueue rxQueue;
  struct rte_mempool* dataMp; ///< mempool for Data seg0
  struct rte_mempool* indirectMp;
  FaceID face;
  uint16_t nPatterns;
  bool wantNackNoRoute; ///< whether to Nack Interests not matching any pattern

  ThreadStopFlag stop;
  uint64_t nNoMatch;
  uint64_t nAllocError;
  pcg32_random_t replyRng;

  PingServerPattern pattern[PINGSERVER_MAX_PATTERNS];
} PingServer;

__attribute__((nonnull)) int
PingServer_Run(PingServer* server);

#endif // NDNDPDK_PINGSERVER_SERVER_H
