# TerraXcel Terraform Provider

Who would like to manage their Excel documents with a tool that is dedicated for infrastrucutre, you may ask. But the more you think about it, the more it make sense, let me explain with a story. 

It is early 2020, covid is not the only thing that is spreading in the world. In every company-wide-meeting you start to hear about this new initative, some call it digital transformation and some call it becoming data-driven, but the bottom line is that hundreds of of hundreds of dollars are going to be spent on new tools, platforms and overpaid engineers to make this happen.

It was a new era, where two main actors are fighting for the throne, the chilly Snowflake and the fiery Databricks. Two great tools with even greater price tags. The merch and swag were being distributed, the dinners and drinks were being served, RFC was approved and the credit cards details was shared.

After all these dollars and hours spent, the company was now ready to start the journey of becoming data-driven. This time no one was building a data warehouse or a data lake, this was so 2010, this time it was all about the Lakehouse. With this new revolutionary architecture(basically just a blob storage and a database) a simple dashboarding tool connected to a table was not enough, this time we needed self-service analytics built upon something called a data catalog. The company was now done, there was nothing more to achive, becoming data driven was the last piece of the puzzle, or so they thought.

Fast forward to the end of the financial year, this year will be nothing like last year, all the groundwork had been laid by the data-team, all the financial data had checked in as a permanent vip-guest in the lakehouse, airflow was continously creating reports and the dashboards were shining brigther than ever. The team is ready to take a early Friday, because thats what you do on Fridays, you know worklife balance, but suddenly you hear the famous _pling_*.

_\*the microsoft teams pling_ 

It's John from finance, "Hey, we are closing the books but the numbers seems to be off, could you send me a updated report?", the team is now confused, what does he mean with a "updated report", the numbers are in the lakehouse, the dashboards are up to date and no errors have been reported, it should be updated, it's streaming analytics, right? The team decides to bite, "the numbers should be updated in the Lakehouse, we have not seen any failed runs the last couple of days, what do you mean with an updated report?", the team wait nervously. The _pling_* comes again.

_\*once again the microsoft teams pling_ 

"Ehm what? I wish I was overpaid like you guys, but I don't have money for a Lakehouse, I am actually in the office, but I would prefer the reports to be digital, maybe you can use the state of the art excel-spreadsheet, I think it is called financial_report_v1_final_2_v2_new_final_final." - John

\*The scene fades out\* 

Since the backbone of you company, the financial reporting anyway going to be done in excel, maybe we were wrong about the whole data-driven thing. Maybe it was the dinner at the fancy restaurant that made us think that we could change the world, maybe the answers was to be found in the old ways, maybe the answer was excel.

If you can't beat them, join them. Let's treat these excel documents as the infrastructure they are, let's use Terraform to manage them. Let's force John make a pull request to change the numbers in the excel document, let's force him to write a commit message, let's make John wait on the DevOps team for review, let's make John's life a living hell. Not because we hate John, but because we love infrastructure and know the importance of version control. Lets do TerraXcel.


## Overview

The TerraXcel Terraform Provider allows you to interact with a TerraXcel server to manage resources, such as creating and manipulating Excel documents (workbooks, sheets, and cells) through Terraform.

## Features

- Create and manage Excel documents.
- Operations on workbooks, sheets, and cells.

## Prerequisites

- A running TerraXcel server.
- Server URL and authentication token.

## Installation

The installation process is manual at the moment. Follow the general guidelines for manual Terraform provider installation, using the `go install` command and other necessary steps.

## Configuration

Configure the provider by specifying the `host` and `token` parameters.

```hcl
provider "terraXcel" {
  host  = "your-server-url"
  token = "your-authentication-token"
}
```

### Parameters

- `host` (Required): URL of the TerraXcel server.
- `token` (Required & Sensitive): Authentication token for the server.

## Usage Example

Here’s an example of how to create a workbook using the TerraXcel provider:

```hcl
resource "terraXcel_workbook" "example" {
  file_name   = "example"
  folder_path = "/path/to/store"
  extension   = "xlsx"
}
```

### Workbook Resource Parameters

- `id` (Computed): Unique ID of the workbook.
- `file_name` (Required): Name of the workbook file.
- `folder_path` (Required): Path where the workbook will be stored.
- `extension` (Required): File extension for the workbook (e.g., "xlsx").
- `last_updated` (Computed): Timestamp of when the workbook was last updated.

## Support and Troubleshooting

For issues, refer to the TerraXcel server documentation and support channels, or consult the broader Terraform community for help.

## Contribution

Contributions are welcome! We encourage the community to contribute to the TerraXcel Terraform Provider project to improve its functionality and performance. Please review the contribution guidelines (if available) before making a contribution.

## License

The TerraXcel Terraform Provider is distributed under a specified license. Review the license documentation accompanying the distribution for more details and ensure compliance with its terms during use and redistribution.

## Contact

For more information, support, and assistance related to the TerraXcel Terraform Provider, please reach out through the official TerraXcel contact channels.
