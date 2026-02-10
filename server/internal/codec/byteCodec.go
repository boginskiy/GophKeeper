package codec

import (
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/infra"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/repo"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
)

type ByteDecoderEncoder struct {
	FileService infra.Filer
	FileHandler utils.FileHandler
	Repo        repo.RepoCreateReader[*model.Bytes]
}

func NewByteDecoderEncoder(
	fileService infra.Filer,
	fileHandler utils.FileHandler,
	repo repo.RepoCreateReader[*model.Bytes]) *ByteDecoderEncoder {

	return &ByteDecoderEncoder{
		FileService: fileService,
		FileHandler: fileHandler,
		Repo:        repo}
}

// DecoderUploadStreamCtx
func (b *ByteDecoderEncoder) DecoderUploadStreamCtx(stream rpc.ByterService_UploadServer) (*model.Bytes, error) {
	modBytes := &model.Bytes{}

	err := modBytes.InsertValuesFromCtx(stream.Context())
	if err != nil {
		return nil, err
	}

	// File for data saving
	file, path, err := b.FileService.CreateFileInStore(modBytes)
	if err != nil {
		return nil, errs.ErrCreateFile.Wrap(err)
	}

	modBytes.Descr, modBytes.Path = file, path
	return modBytes, nil
}

// EncodeUploadBytesResponse
func (a *ByteDecoderEncoder) EncodeUploadBytesResponse(res *model.Bytes) (*rpc.UploadBytesResponse, error) {
	return &rpc.UploadBytesResponse{
		Status:       "uploaded",
		UpdatedAt:    utils.ConvertDtStr(res.UpdatedAt),
		SentSize:     res.SentSize,
		ReceivedSize: res.ReceivedSize,
	}, nil
}

// DecoderUnloadStreamCtx
func (b *ByteDecoderEncoder) DecoderUnloadStreamCtx(stream rpc.ByterService_UnloadServer) (*model.Bytes, error) {
	// Info from context.
	fileName, err := infra.TakeClientValueFromCtx(stream.Context(), "file_name", 0)
	if err != nil {
		return nil, err
	}

	owner, err := infra.TakeServerValStrFromCtx(stream.Context(), infra.EmailCtx)
	if err != nil {
		return nil, err
	}

	// Take record from DataBase.
	modBytes, err := b.Repo.ReadRecord(stream.Context(), &model.Bytes{Name: fileName, Owner: owner})
	if err != nil {
		return nil, err
	}

	return modBytes, nil
}

// EncodeUploadBytesResponse
func (a *ByteDecoderEncoder) EncodeReadBytesResponse(res *model.Bytes) (*rpc.ReadBytesResponse, error) {
	return &rpc.ReadBytesResponse{
		Status:    "read",
		Type:      a.FileHandler.GetTypeFile(res.Name),
		CreatedAt: utils.ConvertDtStr(res.CreatedAt),
	}, nil
}

// EncodeReadAllBytesResponse
func (a *ByteDecoderEncoder) EncodeReadAllBytesResponse(ress []*model.Bytes) (*rpc.ReadAllBytesResponse, error) {
	bytesResponses := make([]*rpc.BytesResponse, len(ress))

	for i, bytes := range ress {
		bytesResponses[i] = &rpc.BytesResponse{
			UpdatedAt: utils.ConvertDtStr(bytes.UpdatedAt),
			TotalSize: bytes.ReceivedSize,
			Name:      bytes.Name,
		}
	}
	return &rpc.ReadAllBytesResponse{
		Status:         "read",
		BytesResponses: bytesResponses}, nil
}

// EncodeDeleteBytesResponse
func (a *ByteDecoderEncoder) EncodeDeleteBytesResponse(status string) (*rpc.DeleteBytesResponse, error) {
	return &rpc.DeleteBytesResponse{
		Status:    status,
		DeletedAt: utils.ConvertDtStr(time.Now())}, nil
}
