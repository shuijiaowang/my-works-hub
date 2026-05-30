import service from '@/utils/request.js'

export const fetchAllProjects = () => {
  return service({
    url: '/projects/all',
    method: 'get',
    donNotShowLoading: true,
  })
}

export const fetchProjectDetail = (folderName) => {
  return service({
    url: `/projects/${encodeURIComponent(folderName)}`,
    method: 'get',
    donNotShowLoading: true,
  })
}

export const updateProject = (folderName, payload) => {
  return service({
    url: `/admin/projects/${encodeURIComponent(folderName)}`,
    method: 'put',
    data: payload,
  })
}

export const fetchProjectMedia = (folderName) => {
  return service({
    url: `/projects/${encodeURIComponent(folderName)}/media`,
    method: 'get',
    donNotShowLoading: true,
  })
}

export const fetchAdminProjectMedia = (folderName) => {
  return service({
    url: `/admin/projects/${encodeURIComponent(folderName)}/media`,
    method: 'get',
    donNotShowLoading: true,
  })
}

export const uploadProjectMedia = (folderName, file) => {
  const form = new FormData()
  form.append('file', file)
  return service({
    url: `/admin/projects/${encodeURIComponent(folderName)}/media`,
    method: 'post',
    data: form,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export const deleteProjectMedia = (folderName, mediaId) => {
  return service({
    url: `/admin/projects/${encodeURIComponent(folderName)}/media/${mediaId}`,
    method: 'delete',
  })
}

export const moveProjectMedia = (folderName, mediaId, direction) => {
  return service({
    url: `/admin/projects/${encodeURIComponent(folderName)}/media/${mediaId}/move`,
    method: 'post',
    data: { direction },
  })
}

export const uploadProjectZip = (folderName, file) => {
  const form = new FormData()
  form.append('file', file)
  return service({
    url: `/admin/projects/${encodeURIComponent(folderName)}/zip`,
    method: 'post',
    data: form,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export const deleteProjectZip = (folderName, fileName) => {
  return service({
    url: `/admin/projects/${encodeURIComponent(folderName)}/zip/${encodeURIComponent(fileName)}`,
    method: 'delete',
  })
}

export const downloadProjectZip = (folderName, fileName) => {
  return service({
    url: `/admin/projects/${encodeURIComponent(folderName)}/zip/${encodeURIComponent(fileName)}`,
    method: 'get',
    responseType: 'blob',
    donNotShowLoading: true,
  })
}
