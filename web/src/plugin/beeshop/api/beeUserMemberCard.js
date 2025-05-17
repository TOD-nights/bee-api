import service from '@/utils/request'


export const userMemberCardPage = (params) => {
  return service({
    url: '/bee-shop/user-member-card/page',
    method: 'get',
    params
  })
}

export const userMemberCardSave = (data) => {
  return service({
    url: '/bee-shop/user-member-card/save',
    method: 'post',
    data
  })
}

export const deleteOneById = (id) => {
  return service({
    url: '/bee-shop/user-member-card/deleteOneById',
    method: 'delete',
    params: { id }
  })
}
export const recoverOneById = (id) => {
  return service({
    url: '/bee-shop/user-member-card/recoverOneById',
    method: 'post',
    params: { id }
  })
}
export const infoById = (id) => {
  return service({
    url: '/bee-shop/user-member-card/info?id='+id,
    method: 'get',
  })
}


