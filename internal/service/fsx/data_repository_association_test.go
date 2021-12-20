package fsx_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/fsx"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tffsx "github.com/hashicorp/terraform-provider-aws/internal/service/fsx"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccFSxDataRepositoryAssociation_basic(t *testing.T) {
	var association fsx.DataRepositoryAssociation
	resourceName := "aws_fsx_data_repository_association.test"
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	bucketPath := fmt.Sprintf("s3://%s", bucketName)
	fileSystemPath := "/test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxDataRepositoryAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFsxDataRepositoryAssociationFileSystemPathConfig(bucketName, fileSystemPath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxDataRepositoryAssociationExists(resourceName, &association),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "fsx", regexp.MustCompile(`association/fs-.+/dra-.+`)),
					resource.TestCheckResourceAttr(resourceName, "batch_import_meta_data_on_create", "false"),
					resource.TestCheckResourceAttr(resourceName, "data_repository_path", bucketPath),
					resource.TestMatchResourceAttr(resourceName, "file_system_id", regexp.MustCompile(`fs-.+`)),
					resource.TestCheckResourceAttr(resourceName, "file_system_path", fileSystemPath),
					resource.TestCheckResourceAttrSet(resourceName, "imported_file_chunk_size"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_data_in_filesystem"},
			},
		},
	})
}

func TestAccFSxDataRepositoryAssociation_s3AutoExportPolicy(t *testing.T) {
	var association fsx.DataRepositoryAssociation
	resourceName := "aws_fsx_data_repository_association.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxDataRepositoryAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFsxDataRepositoryAssociationS3AutoExportPolicyConfig([]string{"NEW", "CHANGED", "DELETED"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxDataRepositoryAssociationExists(resourceName, &association),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_export_policy.events.0", "NEW"),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_export_policy.events.1", "CHANGED"),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_export_policy.events.2", "DELETED"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_data_in_filesystem"},
			},
		},
	})
}

func TestAccFSxDataRepositoryAssociation_s3AutoExportPolicyUpdate(t *testing.T) {
	var association1, association2 fsx.DataRepositoryAssociation
	resourceName := "aws_fsx_data_repository_association.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxDataRepositoryAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFsxDataRepositoryAssociationS3AutoExportPolicyConfig([]string{"NEW", "CHANGED", "DELETED"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxDataRepositoryAssociationExists(resourceName, &association1),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_export_policy.events.0", "NEW"),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_export_policy.events.1", "CHANGED"),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_export_policy.events.2", "DELETED"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_data_in_filesystem"},
			},
			{
				Config: testAccFsxDataRepositoryAssociationS3AutoExportPolicyConfig([]string{"NEW"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxDataRepositoryAssociationExists(resourceName, &association2),
					testAccCheckFsxDataRepositoryAssociationNotRecreated(&association1, &association2),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_export_policy.events.0", "NEW"),
				),
			},
		},
	})
}

func TestAccFSxDataRepositoryAssociation_s3AutoImportPolicy(t *testing.T) {
	var association fsx.DataRepositoryAssociation
	resourceName := "aws_fsx_data_repository_association.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxDataRepositoryAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFsxDataRepositoryAssociationS3AutoImportPolicyConfig([]string{"NEW", "CHANGED", "DELETED"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxDataRepositoryAssociationExists(resourceName, &association),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_import_policy.events.0", "NEW"),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_import_policy.events.1", "CHANGED"),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_import_policy.events.2", "DELETED"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_data_in_filesystem"},
			},
		},
	})
}

func TestAccFSxDataRepositoryAssociation_s3AutoImportPolicyUpdate(t *testing.T) {
	var association1, association2 fsx.DataRepositoryAssociation
	resourceName := "aws_fsx_data_repository_association.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(fsx.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, fsx.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckFsxDataRepositoryAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFsxDataRepositoryAssociationS3AutoImportPolicyConfig([]string{"NEW", "CHANGED", "DELETED"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxDataRepositoryAssociationExists(resourceName, &association1),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_import_policy.events.0", "NEW"),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_import_policy.events.1", "CHANGED"),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_import_policy.events.2", "DELETED"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_data_in_filesystem"},
			},
			{
				Config: testAccFsxDataRepositoryAssociationS3AutoImportPolicyConfig([]string{"NEW"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFsxDataRepositoryAssociationExists(resourceName, &association2),
					testAccCheckFsxDataRepositoryAssociationNotRecreated(&association1, &association2),
					resource.TestCheckResourceAttr(resourceName, "s3.auto_import_policy.events.0", "NEW"),
				),
			},
		},
	})
}

func testAccCheckFsxDataRepositoryAssociationExists(resourceName string, assoc *fsx.DataRepositoryAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).FSxConn

		association, err := tffsx.FindDataRepositoryAssociationByID(conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		if association == nil {
			return fmt.Errorf("FSx Lustre Data Repository Association (%s) not found", rs.Primary.ID)
		}

		*assoc = *association

		return nil
	}
}

func testAccCheckFsxDataRepositoryAssociationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).FSxConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_fsx_lustre_file_system" {
			continue
		}

		filesystem, err := tffsx.FindFileSystemByID(conn, rs.Primary.ID)
		if tfresource.NotFound(err) {
			continue
		}

		if filesystem != nil {
			return fmt.Errorf("FSx Lustre File System (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckFsxDataRepositoryAssociationNotRecreated(i, j *fsx.DataRepositoryAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(i.AssociationId) != aws.StringValue(j.AssociationId) {
			return fmt.Errorf("FSx Data Repository Association (%s) recreated", aws.StringValue(i.AssociationId))
		}

		return nil
	}
}

func testAccCheckFsxDataRepositoryAssociationRecreated(i, j *fsx.DataRepositoryAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(i.AssociationId) == aws.StringValue(j.AssociationId) {
			return fmt.Errorf("FSx Data Repository Association (%s) not recreated", aws.StringValue(i.AssociationId))
		}

		return nil
	}
}

func testAccDataRepositoryAssociationBucketConfig(bucketName string) string {
	return acctest.ConfigCompose(testAccLustreFileSystemBaseConfig(), fmt.Sprintf(`
resource "aws_fsx_lustre_file_system" "test" {
  storage_capacity = 1200
  subnet_ids       = [aws_subnet.test1.id]
  deployment_type  = "PERSISTENT_2"
  per_unit_storage_throughput = 125
}

resource "aws_s3_bucket" "test" {
  acl    = "private"
  bucket = %[1]q
}
`, bucketName))
}

func testAccFsxDataRepositoryAssociationFileSystemPathConfig(bucketName, fileSystemPath string) string {
	return acctest.ConfigCompose(testAccDataRepositoryAssociationBucketConfig(bucketName), fmt.Sprintf(`
resource "aws_fsx_data_repository_association" "test" {
  file_system_id = aws_fsx_lustre_file_system.test.id
  data_repository_path = "s3://%[1]s"
  file_system_path = %[2]q
}
`, bucketName, fileSystemPath))
}

func testAccFsxDataRepositoryAssociationS3AutoExportPolicyConfig(events []string) string {
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	bucketPath := fmt.Sprintf("s3://%s", bucketName)
	fileSystemPath := "/test"
	eventsString := strings.Replace(fmt.Sprintf("%q", events), " ", ", ", -1)
	return acctest.ConfigCompose(testAccDataRepositoryAssociationBucketConfig(bucketName), fmt.Sprintf(`
resource "aws_fsx_data_repository_association" "test" {
  file_system_id = aws_fsx_lustre_file_system.test.id
  data_repository_path = %[1]q
  file_system_path = %[2]q

  s3 {
	  auto_export_policy {
		  events = %[3]s
	  }
  }
}
`, bucketPath, fileSystemPath, eventsString))
}

func testAccFsxDataRepositoryAssociationS3AutoImportPolicyConfig(events []string) string {
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	bucketPath := fmt.Sprintf("s3://%s", bucketName)
	fileSystemPath := "/test"
	eventsString := strings.Replace(fmt.Sprintf("%q", events), " ", ", ", -1)
	return acctest.ConfigCompose(testAccDataRepositoryAssociationBucketConfig(bucketName), fmt.Sprintf(`
resource "aws_fsx_data_repository_association" "test" {
  file_system_id = aws_fsx_lustre_file_system.test.id
  data_repository_path = %[1]q
  file_system_path = %[2]q

  s3 {
	  auto_import_policy {
		  events = %[3]s
	  }
  }
}
`, bucketPath, fileSystemPath, eventsString))
}
