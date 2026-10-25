export default function errorMessage(err, fallback) {
  if (!err.response) {
    return 'Tidak dapat terhubung ke server. Pastikan backend menyala.'
  }
  if (err.response.status === 404 && !err.response.data?.message) {
    return 'Alamat API tidak ditemukan. Restart backend setelah mengubah kode.'
  }
  const message = err.response.data?.message || fallback
  const detail = err.response.data?.error ? ` (${err.response.data.error})` : ''
  return message + detail
}
