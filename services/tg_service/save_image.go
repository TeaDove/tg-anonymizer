package tg_service

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"tg-anonymizer/suppliers/s3_supplier"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func downloadFile(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	return resp.Body, nil
}

func (r *Service) handlePhotoPrivateMessage(
	ctx context.Context,
	wg *sync.WaitGroup,
	update *tgbotapi.Update,
) error {
	defer wg.Done()

	photo := update.Message.Photo[len(update.Message.Photo)-1]

	file, err := r.bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
	if err != nil {
		return errors.Wrap(err, "failed to get file")
	}

	url := file.Link(r.bot.Token)
	closer, err := downloadFile(ctx, url)
	if err != nil {
		return errors.Wrap(err, "failed to download file")
	}

	link, err := r.s3Supplier.PutObject(ctx, &s3_supplier.PutObjectInput{
		Key:         photo.FileID,
		Body:        closer,
		Expires:     time.Now().UTC().Add(10 * 24 * time.Hour),
		ContentType: "image/jpeg",
	})
	if err != nil {
		return errors.Wrap(err, "failed to put object")
	}

	err = r.reply(update, "Link to file: %s", link)
	if err != nil {
		return errors.Wrap(err, "failed to reply")
	}

	return nil
}
