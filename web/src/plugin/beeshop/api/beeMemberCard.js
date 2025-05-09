import service from '@/utils/request'


export const memberCardPage = (params) => {
  return service({
    url: '/bee-shop/member-card/page',
    method: 'get',
    params
  })
}

export const memberCardSave = (data) => {
  return service({
    url: '/bee-shop/member-card/save',
    method: 'post',
    data
  })
}

export const deleteOneById = (id) => {
  return service({
    url: '/bee-shop/member-card/deleteOneById',
    method: 'delete',
    params: { id }
  })
}
export const recoverOneById = (id) => {
  return service({
    url: '/bee-shop/member-card/recoverOneById',
    method: 'post',
    params: { id }
  })
}
export const infoById = (id) => {
  return service({
    url: '/bee-shop/member-card/info?id='+id,
    method: 'get',
  })
}


