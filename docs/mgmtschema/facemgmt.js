(function(exports){
exports.provideDefinitions = function(ctx) {

ctx.declareType('iface.FaceUri', {
  type: 'string',
  format: 'uri',
});

ctx.declareType('iface.Counters', {
  type: 'object',
  properties: {
    RxFrames: ctx.useType('counter'),
    RxOctets: ctx.useType('counter'),
    L2DecodeErrs: ctx.useType('counter'),
    ReassBad: ctx.useType('counter'),
    ReassGood: ctx.useType('counter'),
    L3DecodeErrs: ctx.useType('counter'),
    RxInterests: ctx.useType('counter'),
    RxData: ctx.useType('counter'),
    RxNacks: ctx.useType('counter'),
    FragGood: ctx.useType('counter'),
    FragBad: ctx.useType('counter'),
    TxAllocErrs: ctx.useType('counter'),
    TxQueued: ctx.useType('counter'),
    TxDropped: ctx.useType('counter'),
    TxInterests: ctx.useType('counter'),
    TxData: ctx.useType('counter'),
    TxNacks: ctx.useType('counter'),
    TxFrames: ctx.useType('counter'),
    TxOctets: ctx.useType('counter'),
  },
});

ctx.declareType('dpdk.EthStats', {
  type: 'object',
});

ctx.declareType('socketface.ExCounters', {
  type: 'object',
});

ctx.declareType('facemgmt.FaceInfo', {
  type: 'object',
  properties: {
    Id: ctx.useType('iface.FaceId'),
    LocalUri: ctx.useType('iface.FaceUri'),
    RemoteUri: ctx.useType('iface.FaceUri'),
    IsDown: ctx.useType('boolean'),
    Counters: ctx.useType('iface.Counters'),
    ExCounters: {
      oneOf: [
        ctx.useType('dpdk.EthStats'),
        ctx.useType('socketface.ExCounters'),
        true,
      ],
    },
    Latency: ctx.useType('running_stat.Snapshot'),
  },
});

ctx.declareType('facemgmt.IdArg', ctx.markAllRequired({
  type: 'object',
  properties: {
    Id: ctx.useType('iface.FaceId'),
  },
}));

ctx.declareMethod('Face.List', 'null', 'iface.FaceId[]');

ctx.declareMethod('Face.Get', 'facemgmt.IdArg', 'facemgmt.FaceInfo');

ctx.declareMethod('Face.Create',
  ctx.markAllRequired({
    type: 'object',
    properties: {
      LocalUri: ctx.useType('iface.FaceUri'),
      RemoteUri: ctx.useType('iface.FaceUri'),
    },
  }),
  'facemgmt.IdArg');

ctx.declareMethod('Face.Destroy', 'facemgmt.IdArg', 'null');

};
})(exports);