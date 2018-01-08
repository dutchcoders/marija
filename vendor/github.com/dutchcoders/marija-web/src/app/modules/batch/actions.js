export function batchActions(...actions) {
  console.debug("batchActions", actions);
  return {
    type: 'BATCH_ACTIONS',
    actions: actions
  };
}
