/**
 * 权限控制
 * 基于 JWT 中的 principal_kind 字段
 */
export default function access(initialState: {
  currentUser?: API.CurrentUser;
}) {
  const { currentUser } = initialState;
  const kind = currentUser?.principal_kind;

  return {
    canAdmin: kind === 'admin',
    canAgent: kind === 'agent' || kind === 'admin',
    canService: kind === 'service' || kind === 'admin',
  };
}
