import service from '@/utils/request.js'

export const fetchAllProjects = () => {
  return service({
    url: '/projects/all',
    method: 'get',
    donNotShowLoading: true,
  })
}

export const fetchProjectDetail = (id) => {
  return service({
    url: `/projects/${id}`,
    method: 'get',
    donNotShowLoading: true,
  })
}

export const updateProject = (id, payload) => {
  return service({
    url: `/admin/projects/${id}`,
    method: 'put',
    data: payload,
  })
}

export const fetchProjectMedia = (id) => {
  return service({
    url: `/projects/${id}/media`,
    method: 'get',
    donNotShowLoading: true,
  })
}

export const fetchAdminProjectMedia = (id) => {
  return service({
    url: `/admin/projects/${id}/media`,
    method: 'get',
    donNotShowLoading: true,
  })
}

export const uploadProjectMedia = (id, file) => {
  const form = new FormData()
  form.append('file', file)
  return service({
    url: `/admin/projects/${id}/media`,
    method: 'post',
    data: form,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export const deleteProjectMedia = (id, mediaId) => {
  return service({
    url: `/admin/projects/${id}/media/${mediaId}`,
    method: 'delete',
  })
}

export const moveProjectMedia = (id, mediaId, direction) => {
  return service({
    url: `/admin/projects/${id}/media/${mediaId}/move`,
    method: 'post',
    data: { direction },
  })
}

export const uploadProjectZip = (id, file) => {
  const form = new FormData()
  form.append('file', file)
  return service({
    url: `/admin/projects/${id}/zip`,
    method: 'post',
    data: form,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export const deleteProjectZip = (id, fileName) => {
  return service({
    url: `/admin/projects/${id}/zip/${encodeURIComponent(fileName)}`,
    method: 'delete',
  })
}

export const downloadProjectZip = (id, fileName) => {
  return service({
    url: `/admin/projects/${id}/zip/${encodeURIComponent(fileName)}`,
    method: 'get',
    responseType: 'blob',
    donNotShowLoading: true,
  })
}

