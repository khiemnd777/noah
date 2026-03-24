
export interface StaffModel {
  id: number;
  departmentId?: number | null;
  name: string;
  password?: string;
  email: string;
  phone?: string;
  active?: boolean;
  avatar?: string;
  qrCode?: string;
  roleIds?: number[];
  customFields?: Record<string, any> | null;
  createdAt: string;
  updatedAt: string;
}
