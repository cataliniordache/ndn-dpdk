#ifndef NDN_DPDK_MGMT_HRLOG_WRITER_H
#define NDN_DPDK_MGMT_HRLOG_WRITER_H

/// \file

#include "../../dpdk/eal/thread.h"
#include "post.h"

typedef struct HrlogWriter
{
  ;
} HrlogWriter;

/** \brief Write high resolution logs to a file.
 *  \param nSkip how many initial entries to discard.
 *  \param nTotal how many entries to collect.
 */
int
Hrlog_RunWriter(const char* filename,
                int nSkip,
                int nTotal,
                ThreadStopFlag* stop);

#endif // NDN_DPDK_MGMT_HRLOG_WRITER_H
