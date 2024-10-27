resource "discue_domain" "test_domain" {
  alias    = "my-first-domain"
  hostname = "discue.io"
  port     = 443
}

resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener" {
  queue_id = discue_queue.test_queue.id

  alias        = "my-first-listener"
  liveness_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/live"
  notify_url   = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/notify"
}
