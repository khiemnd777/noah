import { type Perm, usePermissionChecks, useRoleChecks } from "@core/auth/rbac-utils";
import type { MenuItem } from "@core/module/types";

/**
 * Lọc menu theo cả roles và permissions:
 * - Nếu item không khai báo roles/permissions → pass
 * - Nếu có roles: dùng requireAll ? ALL : ANY
 * - Nếu có permissions: dùng requireAll ? ALL : ANY
 * - Cuối cùng AND hai vế (roleOK && permOK)
 */
export function useMenuByAccess(items: MenuItem[]) {
  const { hasAnyRole, hasAllRoles } = useRoleChecks();
  const { hasAnyPermissions, hasAllPermissions } = usePermissionChecks();

  const filterItem = (it: MenuItem): MenuItem | null => {
    const requireAll = !!it.requireAll;

    let roleOK = true;
    if (it.roles?.length) {
      roleOK = requireAll ? hasAllRoles(it.roles) : hasAnyRole(it.roles);
    }

    let permOK = true;
    if (it.permissions?.length) {
      permOK = requireAll
        ? hasAllPermissions(it.permissions as Perm[])
        : hasAnyPermissions(it.permissions as Perm[]);
    }

    if (!roleOK || !permOK) {
      return null;
    }

    const subItems = it.subItems
      ?.map(filterItem)
      .filter((child): child is MenuItem => child !== null);

    if ((subItems?.length ?? 0) > 0) {
      return {
        ...it,
        subItems,
      };
    }

    if (!it.to) {
      return null;
    }

    return {
      ...it,
      subItems: undefined,
    };
  };

  return items
    .map(filterItem)
    .filter((item): item is MenuItem => item !== null);
}

/* Ví dụ sử dụng useMenuByAccess:
const RAW_ITEMS = [
  { key: "home", label: "Dashboard", to: "/" },
  { key: "posts", label: "Posts", to: "/posts", roles: ["editor"], permissions: ["post.manage"] },
  { key: "admin", label: "Admin", to: "/admin", roles: ["admin"], permissions: ["system.manage"] },
];

export function SideMenu() {
  const items = useMenuByAccess(RAW_ITEMS);
  return (
    <nav>
      {items.map((m) => (
        <Link key={m.key} to={m.to}>{m.label}</Link>
      ))}
    </nav>
  );
}
*/
