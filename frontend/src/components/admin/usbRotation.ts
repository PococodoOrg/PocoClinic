export interface UsbDrive {
  id: string;
  label: string;
  name: string;
  weekdays: number[];
}

const STORAGE_KEY = 'poco-usb-rotation';

export const DEFAULT_USB_DRIVES: UsbDrive[] = [
  { id: 'mon', label: 'MON', name: 'Monday drive', weekdays: [1] },
  { id: 'wed', label: 'WED', name: 'Wednesday drive', weekdays: [3] },
  { id: 'fri', label: 'FRI', name: 'Friday drive', weekdays: [5] },
];

export function loadUsbDrives(): UsbDrive[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return DEFAULT_USB_DRIVES;
    }
    const parsed = JSON.parse(raw) as UsbDrive[];
    if (!Array.isArray(parsed) || parsed.length === 0) {
      return DEFAULT_USB_DRIVES;
    }
    return parsed;
  } catch {
    return DEFAULT_USB_DRIVES;
  }
}

export function saveUsbDrives(drives: UsbDrive[]): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(drives));
}

export function todaysDrive(drives: UsbDrive[], date = new Date()): UsbDrive | undefined {
  const weekday = date.getDay();
  return drives.find((drive) => drive.weekdays.includes(weekday));
}

export function weekdayLabel(weekdays: number[]): string {
  const names = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  return weekdays.map((day) => names[day]).join(', ');
}
