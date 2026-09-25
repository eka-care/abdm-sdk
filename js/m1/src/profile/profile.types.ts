export interface AbhaProfile {
  ABHANumber: string;
  firstName: string;
  middleName?: string;
  lastName?: string;
  dob?: string;
  gender: string;
  mobile?: string;
  email?: string;
  phrAddress?: string[];
  address?: string;
  districtCode?: string;
  districtName?: string;
  stateCode?: string;
  stateName?: string;
  pinCode?: string;
  abhaType?: string;
  abhaStatus?: string;
  photo?: string;
  [key: string]: unknown;
}
