# Generated by Django 2.2.2 on 2019-06-30 06:46

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('users', '0008_auto_20190630_0646'),
    ]

    operations = [
        migrations.AlterField(
            model_name='profile',
            name='password',
            field=models.CharField(default='blockchain', max_length=20),
        ),
    ]
