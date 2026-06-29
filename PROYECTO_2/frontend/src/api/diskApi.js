import axiosClient from "./axiosClient";
//diskApi.js: contiene las funciones para manejar las operaciones de disco y partición, incluyendo la creación, eliminación, redimensionamiento, montaje y formateo de discos y particiones, así como la obtención de información sobre discos y particiones montadas. Todas las funciones se comunican con el backend a través de axiosClient.

function backendError(error) {
  return (
    error?.response?.data?.error ||
    error?.response?.data?.message ||
    error?.message ||
    "Error desconocido"
  );
}

async function safeRequest(callback) {
  try {
    return await callback();
  } catch (error) {
    throw new Error(backendError(error));
  }
}

function normalizeSize(size, unit) {
  const numericSize = Number(size);
  const normalizedUnit = unit || "M";

  if (!Number.isFinite(numericSize) || numericSize <= 0) {
    return {
      size: 0,
      unit: normalizedUnit,
    };
  }

  if (Number.isInteger(numericSize)) {
    return {
      size: numericSize,
      unit: normalizedUnit,
    };
  }

  if (normalizedUnit === "M") {
    return {
      size: Math.round(numericSize * 1024 * 1024),
      unit: "B",
    };
  }

  if (normalizedUnit === "K") {
    return {
      size: Math.round(numericSize * 1024),
      unit: "B",
    };
  }

  return {
    size: Math.round(numericSize),
    unit: "B",
  };
}

function normalizeAdd(add, unit) {
  const numericAdd = Number(add);
  const normalizedUnit = unit || "M";

  if (!Number.isFinite(numericAdd) || numericAdd === 0) {
    return {
      add: 0,
      unit: normalizedUnit,
    };
  }

  if (Number.isInteger(numericAdd)) {
    return {
      add: numericAdd,
      unit: normalizedUnit,
    };
  }

  if (normalizedUnit === "M") {
    return {
      add: Math.round(numericAdd * 1024 * 1024),
      unit: "B",
    };
  }

  if (normalizedUnit === "K") {
    return {
      add: Math.round(numericAdd * 1024),
      unit: "B",
    };
  }

  return {
    add: Math.round(numericAdd),
    unit: "B",
  };
}

export async function getDisks() {
  return safeRequest(async () => {
    const response = await axiosClient.get("/disks");
    return response.data.disks || response.data.data || [];
  });
}

export async function createDisk(payload) {
  return safeRequest(async () => {
    const normalized = normalizeSize(payload.size, payload.unit);

    const response = await axiosClient.post("/disks", {
      name: payload.name,
      path: payload.path,
      size: normalized.size,
      unit: normalized.unit,
      fit: payload.fit,
    });

    return response.data.disk || response.data.data;
  });
}

export async function deleteDisk(path) {
  return safeRequest(async () => {
    const response = await axiosClient.post("/disks/delete", { path });
    return response.data;
  });
}

export async function removeDisk(path) {
  return deleteDisk(path);
}

export async function getPartitions(diskPath) {
  return safeRequest(async () => {
    const response = await axiosClient.get("/partitions/list", {
      params: { diskPath },
    });

    return response.data.partitions || response.data.data || [];
  });
}

export async function createPartition(payload) {
  return safeRequest(async () => {
    const normalized = normalizeSize(payload.size, payload.unit);

    const response = await axiosClient.post("/partitions/create", {
      diskPath: payload.diskPath || payload.path,
      path: payload.path || payload.diskPath,
      name: payload.name,
      size: normalized.size,
      unit: normalized.unit,
      type: payload.type,
      fit: payload.fit,
    });

    return response.data.partition || response.data.data;
  });
}

export async function deletePartition(payload) {
  return safeRequest(async () => {
    const response = await axiosClient.post("/partitions/delete", {
      diskPath: payload.diskPath || payload.path,
      path: payload.path || payload.diskPath,
      name: payload.name,
      deleteType: payload.deleteType || payload.delete || "fast",
      delete: payload.delete || payload.deleteType || "fast",
    });

    return response.data;
  });
}

export async function removePartition(payload) {
  return deletePartition(payload);
}

export async function resizePartition(payload) {
  return safeRequest(async () => {
    const normalized = normalizeAdd(payload.add, payload.unit);

    const response = await axiosClient.post("/partitions/resize", {
      diskPath: payload.diskPath || payload.path,
      path: payload.path || payload.diskPath,
      name: payload.name,
      add: normalized.add,
      unit: normalized.unit,
    });

    return response.data.partition || response.data.data;
  });
}

export async function mountPartition(payload) {
  return safeRequest(async () => {
    const response = await axiosClient.post("/partitions/mount", {
      diskPath: payload.diskPath || payload.path,
      path: payload.path || payload.diskPath,
      name: payload.name,
    });

    return response.data.mounted || response.data.data;
  });
}

export async function unmountPartition(id) {
  return safeRequest(async () => {
    const response = await axiosClient.post("/partitions/unmount", { id });
    return response.data;
  });
}

export async function formatPartition(payload) {
  return safeRequest(async () => {
    const response = await axiosClient.post("/mkfs", {
      id: payload.id,
      type: payload.type || "full",
    });

    return response.data.partition || response.data.data;
  });
}

export async function getMountedPartitions() {
  return safeRequest(async () => {
    const response = await axiosClient.get("/mounted");
    return response.data.mounted || response.data.data || [];
  });
}
